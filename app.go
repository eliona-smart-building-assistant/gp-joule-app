//  This file is part of the eliona project.
//  Copyright © 2022 LEICOM iTEC AG. All Rights Reserved.
//  ______ _ _
// |  ____| (_)
// | |__  | |_  ___  _ __   __ _
// |  __| | | |/ _ \| '_ \ / _` |
// | |____| | | (_) | | | | (_| |
// |______|_|_|\___/|_| |_|\__,_|
//
//  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING
//  BUT NOT LIMITED  TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
//  NON INFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
//  DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	"context"
	"gp-joule/apiserver"
	"gp-joule/apiservices"
	"gp-joule/conf"
	"gp-joule/gp_joule"
	"gp-joule/model"
	"net/http"
	"sync"
	"time"

	"github.com/eliona-smart-building-assistant/go-eliona/app"
	"github.com/eliona-smart-building-assistant/go-eliona/asset"
	"github.com/eliona-smart-building-assistant/go-eliona/dashboard"
	"github.com/eliona-smart-building-assistant/go-eliona/frontend"
	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/eliona-smart-building-assistant/go-utils/db"
	utilshttp "github.com/eliona-smart-building-assistant/go-utils/http"
	"github.com/eliona-smart-building-assistant/go-utils/log"
)

func initialization() {
	ctx := context.Background()

	// Necessary to close used init resources
	conn := db.NewInitConnectionWithContextAndApplicationName(ctx, app.AppName())
	defer conn.Close(ctx)

	// Init the app before the first run.
	app.Init(conn, app.AppName(),
		app.ExecSqlFile("conf/init.sql"),
		asset.InitAssetTypeFiles("resources/asset-types/*.json"),
		dashboard.InitWidgetTypeFiles("resources/widget-types/*.json"),
	)
}

var once sync.Once

func collectData() {
	configs, err := conf.GetConfigs(context.Background())
	if err != nil {
		log.Fatal("conf", "Couldn't read configs from DB: %v", err)
		return
	}
	if len(configs) == 0 {
		once.Do(func() {
			log.Info("conf", "No configs in DB. Please configure the app in Eliona.")
		})
		return
	}

	for _, config := range configs {
		if !conf.IsConfigEnabled(config) {
			if conf.IsConfigActive(config) {
				_, _ = conf.SetConfigActiveState(context.Background(), config, false)
			}
			continue
		}

		if !conf.IsConfigActive(config) {
			_, _ = conf.SetConfigActiveState(context.Background(), config, true)
			log.Info("conf", "Collecting initialized with Configuration %d:\n"+
				"Enable: %t\n"+
				"Refresh Interval: %d\n"+
				"Request Timeout: %d\n"+
				"Project IDs: %v\n",
				*config.Id,
				*config.Enable,
				config.RefreshInterval,
				*config.RequestTimeout,
				*config.ProjectIDs)
		}

		common.RunOnceWithParam(func(config apiserver.Configuration) {
			log.Info("main", "Collecting %d started.", *config.Id)
			if err := collectResources(&config); err != nil {
				return // Error is handled in the method itself.
			}
			log.Info("main", "Collecting %d finished.", *config.Id)

			time.Sleep(time.Second * time.Duration(config.RefreshInterval))
		}, config, *config.Id)
	}
}

func collectResources(config *apiserver.Configuration) error {

	// check if project ids are defined, warn if not
	if config.ProjectIDs == nil || len(*config.ProjectIDs) == 0 {
		log.Warn("api", "No project IDs defined in config %d", *config.Id)
		return nil
	}

	// get all clusters from GP Joule API
	clusters, err := gp_joule.GetClusters(config)
	if err != nil {
		log.Error("api", "Error collecting clusters: %v", err)
		return err
	}
	log.Trace("api", "Clusters: %v", clusters)

	// Create asset tree for each project id
	for _, projectId := range *config.ProjectIDs {

		// Create asset tree
		root := model.Root{
			Clusters: clusters,
		}

		// This is ugly but there is no other way to put the config into each element
		root.Config = config
		for _, cluster := range root.Clusters {
			cluster.Config = config
			for _, chargingPoint := range cluster.ChargePoints {
				chargingPoint.Config = config
				for _, connector := range chargingPoint.Connectors {
					connector.Config = config
				}
			}
		}

		// create assets
		assets, err := asset.CreateAssetsAndUpsertData(&root, projectId, nil, nil)

		if err != nil {
			log.Error("eliona", "Error creating assets: %v", err)
			return err
		}

		log.Trace("eliona", "Assets created: %v", assets)
	}

	return nil
}

// listenApi starts the API server and listen for requests
func listenApi() {
	err := http.ListenAndServe(":"+common.Getenv("API_SERVER_PORT", "3000"),
		frontend.NewEnvironmentHandler(
			utilshttp.NewCORSEnabledHandler(
				apiserver.NewRouter(
					apiserver.NewConfigurationAPIController(apiservices.NewConfigurationAPIService()),
					apiserver.NewVersionAPIController(apiservices.NewVersionAPIService()),
					apiserver.NewCustomizationAPIController(apiservices.NewCustomizationAPIService()),
				))))
	log.Fatal("main", "API server: %v", err)
}
