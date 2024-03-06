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

package eliona

import (
	"context"
	"fmt"
	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v2"
	"github.com/eliona-smart-building-assistant/go-eliona/asset"
	"github.com/eliona-smart-building-assistant/go-eliona/client"
	"github.com/eliona-smart-building-assistant/go-utils/log"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"gp-joule/apiserver"
	"gp-joule/appdb"
)

// InitAssets initializes the assets created before. This contains creation of pipeline aggregation and rules for alarms
func InitAssets(config *apiserver.Configuration) error {
	dbAssets, err := appdb.Assets(
		appdb.AssetWhere.ConfigurationID.EQ(*config.Id),
		appdb.AssetWhere.InitVersion.LTE(1),
	).AllG(context.Background())
	if err != nil {
		return err
	}
	for _, dbAsset := range dbAssets {
		err = initAsset(dbAsset)
		if err != nil {
			return err
		}
	}
	return nil
}

func initAsset(dbAsset *appdb.Asset) error {
	if dbAsset == nil {
		return nil
	}
	if dbAsset.InitVersion <= 0 {
		err := initAssetV1(dbAsset)
		if err != nil {
			return err
		}
		dbAsset.InitVersion = 1
		_, err = dbAsset.UpdateG(context.Background(), boil.Whitelist(appdb.AssetColumns.InitVersion))
		if err != nil {
			return err
		}
	}
	if dbAsset.InitVersion <= 1 {
		// Place for init during a patch of new app version
	}
	return nil
}

func initAssetV1(dbAsset *appdb.Asset) error {

	// check if asset still exists in Eliona
	exists, err := asset.ExistAsset(dbAsset.AssetID.Int32)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}

	if dbAsset.AssetType.String == "gp_joule_charge_point" {

		// Set alarm rules
		// ...

	}

	return nil
}

func notifyUser(userId string, projectId string, assetsCreated int) error {
	receipt, _, err := client.NewClient().CommunicationAPI.
		PostNotification(client.AuthenticationContext()).
		Notification(
			api.Notification{
				User:      userId,
				ProjectId: *api.NewNullableString(&projectId),
				Message: *api.NewNullableTranslation(&api.Translation{
					De: api.PtrString(fmt.Sprintf("GP Joule App hat %d neue Assets angelegt. Diese sind nun im Asset-Management verfügbar.", assetsCreated)),
					En: api.PtrString(fmt.Sprintf("GP Joule app added %v new assets. They are now available in Asset Management.", assetsCreated)),
				}),
			}).
		Execute()
	log.Debug("eliona", "posted notification about CAC: %v", receipt)
	if err != nil {
		return fmt.Errorf("posting CAC notification: %v", err)
	}
	return nil
}
