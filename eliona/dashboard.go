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
	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v2"
	"github.com/eliona-smart-building-assistant/go-eliona/asset"
	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/eliona-smart-building-assistant/go-utils/log"
	"gp-joule/conf"
)

func GpJouleDashboard(projectId string) (api.Dashboard, error) {

	dashboard := api.Dashboard{}
	dashboard.Name = "GP Joule"
	dashboard.ProjectId = projectId
	dashboard.Widgets = []api.Widget{}

	dbConnectorAssets, err := conf.GetConnectorsPerProject(context.Background(), projectId)
	if err != nil {
		log.Error("eliona", "Error getting connectors: %v", err)
		return dashboard, err
	}

	for _, dbConnectorAsset := range dbConnectorAssets {

		// check if asset still exists in Eliona
		exists, err := asset.ExistAsset(dbConnectorAsset.AssetID.Int32)
		if err != nil {
			log.Error("eliona", "Error checking asset exists: %v", err)
			return dashboard, err
		}

		if exists {

			// Get sessions asset for this
			dbSessionsLogAsset, err := conf.GetSessionsLog(context.Background(), dbConnectorAsset.ProviderID)
			if err != nil {
				log.Error("eliona", "Error getting sessions log : %v", err)
				return dashboard, err
			}

			// Get charge point asset for this
			dbChargePointAsset, err := conf.GetChargePoint(context.Background(), dbConnectorAsset.ParentProviderID)
			if err != nil {
				log.Error("eliona", "Error getting sessions log : %v", err)
				return dashboard, err
			}

			if dbSessionsLogAsset != nil && dbChargePointAsset != nil {

				widget := api.Widget{
					WidgetTypeName: "gp_joule_connector",
					AssetId:        *api.NewNullableInt32(common.Ptr(dbConnectorAsset.AssetID.Int32)),
					Sequence:       *api.NewNullableInt32(common.Ptr(int32(0))),
					Data: []api.WidgetData{
						{
							ElementSequence: *api.NewNullableInt32(common.Ptr(int32(1))),
							AssetId:         *api.NewNullableInt32(common.Ptr(dbConnectorAsset.AssetID.Int32)),
							Data: map[string]interface{}{
								"aggregatedDataType": "heap",
								"attribute":          "occupied",
								"description":        "Status",
								"key":                "",
								"seq":                0,
								"subtype":            "status",
							},
						},
						{
							ElementSequence: *api.NewNullableInt32(common.Ptr(int32(2))),
							AssetId:         *api.NewNullableInt32(common.Ptr(dbChargePointAsset.AssetID.Int32)),
							Data: map[string]interface{}{
								"aggregatedDataType": "heap",
								"attribute":          "model",
								"description":        "Model",
								"key":                "",
								"seq":                0,
								"subtype":            "info",
							},
						},
						{
							ElementSequence: *api.NewNullableInt32(common.Ptr(int32(2))),
							AssetId:         *api.NewNullableInt32(common.Ptr(dbChargePointAsset.AssetID.Int32)),
							Data: map[string]interface{}{
								"aggregatedDataType": "heap",
								"attribute":          "manufacturer",
								"description":        "Manufacturer",
								"key":                "",
								"seq":                1,
								"subtype":            "info",
							},
						},
						{
							ElementSequence: *api.NewNullableInt32(common.Ptr(int32(2))),
							AssetId:         *api.NewNullableInt32(common.Ptr(dbConnectorAsset.AssetID.Int32)),
							Data: map[string]interface{}{
								"aggregatedDataType": "heap",
								"attribute":          "max_power",
								"description":        "Power",
								"key":                "",
								"seq":                2,
								"subtype":            "info",
							},
						},
						{
							ElementSequence: *api.NewNullableInt32(common.Ptr(int32(2))),
							AssetId:         *api.NewNullableInt32(common.Ptr(dbSessionsLogAsset.AssetID.Int32)),
							Data: map[string]interface{}{
								"aggregatedDataField":  "sum",
								"aggregatedDataRaster": "DECADE",
								"aggregatedDataType":   "pipeline",
								"attribute":            "energy",
								"description":          "Energy total",
								"key":                  "",
								"seq":                  3,
								"subtype":              "input",
							},
						},
						{
							ElementSequence: *api.NewNullableInt32(common.Ptr(int32(2))),
							AssetId:         *api.NewNullableInt32(common.Ptr(dbConnectorAsset.AssetID.Int32)),
							Data: map[string]interface{}{
								"aggregatedDataType": "heap",
								"attribute":          "error_message",
								"description":        "Error message",
								"key":                "",
								"seq":                4,
								"subtype":            "status",
							},
						},
						{
							ElementSequence: *api.NewNullableInt32(common.Ptr(int32(3))),
							AssetId:         *api.NewNullableInt32(common.Ptr(dbSessionsLogAsset.AssetID.Int32)),
							Data: map[string]interface{}{
								"aggregatedDataField":  "sum",
								"aggregatedDataRaster": "DAY",
								"aggregatedDataType":   "pipeline",
								"attribute":            "energy",
								"description":          "Energy",
								"key":                  "",
								"seq":                  0,
								"subtype":              "input",
							},
						},
						{
							ElementSequence: *api.NewNullableInt32(common.Ptr(int32(3))),
							AssetId:         *api.NewNullableInt32(common.Ptr(dbSessionsLogAsset.AssetID.Int32)),
							Data: map[string]interface{}{
								"aggregatedDataField":  "sum",
								"aggregatedDataRaster": "DAY",
								"aggregatedDataType":   "pipeline",
								"attribute":            "duration",
								"description":          "Duration",
								"key":                  "",
								"seq":                  1,
								"subtype":              "input",
							},
						},
					},
					Details: map[string]any{
						"1": map[string]any{
							"tilesConfig": []map[string]any{
								{
									"defaultColorIndex": 7,
									"valueMapping": [][]string{
										{
											"0",
											"Available",
											"#007305",
										},
										{
											"1",
											"Occupied",
											"#9E003D",
										},
									},
								},
							},
						},
						"3": map[string]any{
							"colors": []string{
								"#35c7d5",
								"#656565",
							},
						},
						"size":     1,
						"timespan": 30,
					},
				}

				// add station widget to dashboard
				dashboard.Widgets = append(dashboard.Widgets, widget)
			}
		}
	}

	return dashboard, nil
}
