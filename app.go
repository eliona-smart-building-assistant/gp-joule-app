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
	"fmt"
	"gp-joule/apiserver"
	"gp-joule/apiservices"
	"gp-joule/appdb"
	"gp-joule/conf"
	"gp-joule/eliona"
	"gp-joule/gp_joule"
	"gp-joule/model"
	"math"
	"net/http"
	"sync"
	"time"

	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v2"
	"github.com/volatiletech/sqlboiler/v4/boil"

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

			log.Info("main", "Collecting for config %d started.", *config.Id)
			if err := collectResources(&config); err != nil {
				return // ErrorNotification is handled in the method itself.
			}
			if err := sendSessions(&config); err != nil {
				return // ErrorNotification is handled in the method itself.
			}
			if err := sendErrors(&config); err != nil {
				return // ErrorNotification is handled in the method itself.
			}
			log.Info("main", "Collecting for config %d finished.", *config.Id)

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
		log.Error("api", "ErrorNotification collecting clusters: %v", err)
		return err
	}
	log.Trace("api", "Clusters: %v", clusters)

	// Create asset tree for each project id
	for _, projectId := range *config.ProjectIDs {

		log.Debug("eliona", "Start creating assets for config %d", *config.Id)

		// Create asset tree
		root := model.Root{
			Config:   config,
			Clusters: clusters,
		}

		// create assets
		count, err := asset.CreateAssetsAndUpsertData(&root, projectId, nil, nil)
		if err != nil {
			log.Error("eliona", "ErrorNotification creating assets for config %d: %v", *config.Id, err)
			return err
		}

		log.Debug("eliona", "%d assets created for config %d", count, *config.Id)

		// send notification
		if count > 0 {
			err = eliona.NotifyUser(config.UserId, projectId, &api.Translation{
				De: api.PtrString(fmt.Sprintf("GP Joule App hat %d neue Assets angelegt. Diese sind nun im Asset-Management verfügbar.", count)),
				En: api.PtrString(fmt.Sprintf("GP Joule app added %v new assets. They are now available in Asset Management.", count)),
			})
			if err != nil {
				log.Error("collect", "ErrorNotification notifying users about asset creation: %v", err)
			}
		}
	}

	// init assets
	log.Debug("eliona", "Start init assets for config %d", *config.Id)

	err = eliona.InitAssets(config)
	if err != nil {
		log.Error("eliona", "ErrorNotification creating assets: %v", err)
		return err
	}

	log.Debug("eliona", "Finished init assets for config %d", *config.Id)

	return nil
}

func sendSessions(config *apiserver.Configuration) error {

	dbConnectorAssets, err := conf.GetConnectors(context.Background(), config)
	if err != nil {
		log.Error("eliona", "Error getting connectors: %v", err)
		return err
	}

	log.Debug("eliona", "Start sending sessions for config %d", *config.Id)
	for _, dbConnectorAsset := range dbConnectorAssets {
		var count = 0

		// check if asset still exists in Eliona
		exists, err := asset.ExistAsset(dbConnectorAsset.AssetID.Int32)
		if err != nil {
			log.Error("eliona", "Error checking asset exists: %v", err)
			return err
		}

		if exists {

			// get all sessions
			completedSessions, err := gp_joule.GetCompletedSessions(config, dbConnectorAsset)
			if err != nil {
				log.Error("api", "Error collecting completed sessions: %v", err)
				return err
			}

			// Get sessions asset for this
			dbSessionsLogAsset, err := conf.GetSessionsLog(context.Background(), dbConnectorAsset.ProviderID)
			if err != nil {
				log.Error("eliona", "Error getting sessions log : %v", err)
				return err
			}

			if dbSessionsLogAsset != nil {

				// send sessions to Eliona
				for _, completedSession := range completedSessions {

					// send new session as data to Eliona
					err = asset.UpsertData(api.Data{
						AssetId:   dbSessionsLogAsset.AssetID.Int32,
						Subtype:   "input",
						Timestamp: *api.NewNullableTime(completedSession.SessionEnd),
						Data: map[string]any{
							"count":    1, // Helper attribute to calculate number of charging sessions in Eliona.
							"energy":   int(math.Max(float64(completedSession.MeterTotal), 0)),
							"duration": completedSession.Duration,
						},
					})
					if err != nil {
						log.Error("api", "Error upserting data in Eliona: %v", err)
						return err
					}

					// remember latest timestamp
					dbConnectorAsset.LatestSessionTS = *completedSession.SessionEnd
					_, err = dbConnectorAsset.UpdateG(context.Background(), boil.Whitelist(appdb.AssetColumns.LatestSessionTS))
					if err != nil {
						log.Error("api", "ErrorNotification updating asset latest session timestamp: %v", err)
						return err
					}

					count++
				}
			}
		}

		log.Debug("eliona", "Finished sending %d sessions for asset %d for config %d", count, dbConnectorAsset.AssetID.Int32, *config.Id)
	}

	log.Debug("eliona", "Finished sending sessions for config %d", *config.Id)

	return nil
}

func sendErrors(config *apiserver.Configuration) error {

	dbConnectorAssets, err := conf.GetConnectors(context.Background(), config)
	if err != nil {
		log.Error("eliona", "Error getting connectors: %v", err)
		return err
	}

	log.Debug("eliona", "Start sending errors for config %d", *config.Id)
	for _, dbConnectorAsset := range dbConnectorAssets {
		var openCount = 0
		var resolvedCount = 0

		// check if asset still exists in Eliona
		exists, err := asset.ExistAsset(dbConnectorAsset.AssetID.Int32)
		if err != nil {
			log.Error("eliona", "Error checking asset exists: %v", err)
			return err
		}

		if exists {

			// get all open errors
			errorNotifications, err := gp_joule.GetErrorNotifications(config, dbConnectorAsset)
			if err != nil {
				log.Error("api", "Error collecting error notifications: %v", err)
				return err
			}

			// send all new resolved errors
			for _, errorNotification := range errorNotifications {

				// Close all errors that are resolved until an open error is found
				if errorNotification.ResolvedAt != nil {
					resolvedCount++

					// reset resoled error as data to Eliona
					err = asset.UpsertData(api.Data{
						AssetId:   dbConnectorAsset.AssetID.Int32,
						Subtype:   "status",
						Timestamp: *api.NewNullableTime(errorNotification.ResolvedAt),
						Data: map[string]any{
							"error":         0,
							"error_message": "-",
						},
					})
					if err != nil {
						log.Error("api", "Error upserting data in Eliona: %v", err)
						return err
					}
					// remember latest timestamp
					dbConnectorAsset.LatestErrorTS = *errorNotification.OccurredAt
					_, err = dbConnectorAsset.UpdateG(context.Background(), boil.Whitelist(appdb.AssetColumns.LatestErrorTS))
					if err != nil {
						log.Error("api", "ErrorNotification updating asset latest session timestamp: %v", err)
						return err
					}

				}
			}

			// get all open errors
			errorNotifications, err = gp_joule.GetErrorNotifications(config, dbConnectorAsset)
			if err != nil {
				log.Error("api", "Error collecting error notifications: %v", err)
				return err
			}

			// send all open errors
			for _, errorNotification := range errorNotifications {

				// Send open errors to Eliona
				if errorNotification.ResolvedAt == nil {
					openCount++

					// send new error as data to Eliona
					err = asset.UpsertData(api.Data{
						AssetId:   dbConnectorAsset.AssetID.Int32,
						Subtype:   "status",
						Timestamp: *api.NewNullableTime(errorNotification.OccurredAt),
						Data: map[string]any{
							"error":         openCount,
							"error_message": fmt.Sprintf("%s: %s (%s)", errorNotification.ErrorCode, errorNotification.ErrorInfo, errorNotification.Id),
						},
					})
					if err != nil {
						log.Error("api", "Error upserting data in Eliona: %v", err)
						return err
					}

				}
			}

		}

		log.Debug("eliona", "Finished opening %d new errors for asset %d for config %d", openCount, dbConnectorAsset.AssetID.Int32, *config.Id)
		log.Debug("eliona", "Finished closing %d resolved errors for asset %d for config %d", resolvedCount, dbConnectorAsset.AssetID.Int32, *config.Id)
	}

	log.Debug("eliona", "Finished sending errors for config %d", *config.Id)

	return nil
}

// listenApi starts the API server and listen for requests
func listenApi() {
	log.Info("main", "Starting API server")
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
