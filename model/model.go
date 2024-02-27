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

package model

import (
	"context"
	"fmt"
	"gp-joule/apiserver"
	"gp-joule/conf"
	"time"

	"github.com/eliona-smart-building-assistant/go-eliona/asset"
	"github.com/eliona-smart-building-assistant/go-eliona/utils"
	"github.com/eliona-smart-building-assistant/go-utils/common"
)

// ROOT

type Root struct {
	DefaultImpl
	Clusters []Cluster
}

func (r *Root) GetName() string {
	return "GP Joule"
}

func (r *Root) GetDescription() string {
	return "GP Joule Root"
}

func (r *Root) GetAssetType() string {
	return "gp_joule_root"
}

func (r *Root) GetGAI() string {
	return r.GetAssetType()
}

func (r *Root) GetProviderId() string {
	return ""
}

func (r *Root) GetLocationalChildren() []asset.LocationalNode {
	locationalChildren := make([]asset.LocationalNode, 0)
	for _, cluster := range r.Clusters {
		locationalChildren = append(locationalChildren, common.Ptr(cluster))
	}
	return locationalChildren
}

// CLUSTER

type Cluster struct {
	DefaultImpl
	Name         string        `json:"name" eliona:"name,filterable"`
	ChargePoints []ChargePoint `json:"chargepoints"`
}

func (c *Cluster) GetName() string {
	return c.Name
}

func (c *Cluster) GetDescription() string {
	return ""
}

func (c *Cluster) GetAssetType() string {
	return "gp_joule_cluster"
}

func (c *Cluster) GetGAI() string {
	return c.GetAssetType() + "_" + c.Name
}

func (c *Cluster) GetProviderId() string {
	return c.Name
}

func (c *Cluster) GetLocationalChildren() []asset.LocationalNode {
	locationalChildren := make([]asset.LocationalNode, 0)
	for _, chargingPoint := range c.ChargePoints {
		locationalChildren = append(locationalChildren, common.Ptr(chargingPoint))
	}
	return locationalChildren
}

// CHARGE POINT

type ChargePoint struct {
	DefaultImpl
	ChargePointId       string      `json:"chargepoint_id" eliona:"id,filterable"`
	ChargePointOcppId   string      `json:"chargepoint_ocpp_id"`
	Name                string      `json:"name" eliona:"name,filterable"`
	NameInternal        string      `json:"name_internal" eliona:"name_internal,filterable"`
	Status              string      `json:"status" eliona:"status" subtype:"status"`
	CommunicationStatus int         `json:"communication_status"`
	ConnectorsTotal     int         `json:"connectors_total"`
	ConnectorsFree      int         `json:"connectors_free"`
	ConnectorsFaulted   int         `json:"connectors_faulted"`
	ConnectorsOccupied  int         `json:"connectors_occupied"`
	Manufacturer        string      `json:"manufacturer"`
	Model               string      `json:"model" eliona:"model,filterable" subtype:"info"`
	Lat                 float64     `json:"lat"`
	Long                float64     `json:"long"`
	Street              string      `json:"street"`
	Zip                 int         `json:"zip"`
	City                string      `json:"city"`
	CountryCode         interface{} `json:"country_code"`
	Country             string      `json:"country"`
	Error               interface{} `json:"error"`
	Connectors          []Connector `json:"connectors"`
}

func (cp *ChargePoint) GetName() string {
	return cp.Name
}

func (cp *ChargePoint) GetDescription() string {
	return cp.NameInternal
}

func (cp *ChargePoint) GetAssetType() string {
	return "gp_joule_charge_point"
}

func (cp *ChargePoint) GetGAI() string {
	return cp.GetAssetType() + "_" + cp.ChargePointId
}

func (cp *ChargePoint) GetProviderId() string {
	return cp.ChargePointId
}

func (cp *ChargePoint) GetLocationalChildren() []asset.LocationalNode {
	locationalChildren := make([]asset.LocationalNode, 0)
	for _, connector := range cp.Connectors {
		locationalChildren = append(locationalChildren, common.Ptr(connector))
	}
	return locationalChildren
}

// CONNECTOR

type Connector struct {
	DefaultImpl
	Uuid            string `json:"uuid"`
	EvseId          string `json:"evseid"`
	Status          string `json:"status" eliona:"status" subtype:"status"`
	MaxPower        int    `json:"max_power" eliona:"max_power" subtype:"info"`
	ChargePointType string `json:"chargepoint_type" eliona:"chargepoint_type,filterable" subtype:"info"`
	PlugType        string `json:"plug_type" eliona:"plug_type,filterable" subtype:"info"`
}

func (c *Connector) GetName() string {
	return c.EvseId
}

func (c *Connector) GetDescription() string {
	return c.ChargePointType + " " + c.PlugType
}

func (c *Connector) GetAssetType() string {
	return "gp_joule_connector"
}

func (c *Connector) GetGAI() string {
	return c.GetAssetType() + "_" + c.Uuid
}

func (c *Connector) GetProviderId() string {
	return c.Uuid
}

// SESSION

type ChargingSession struct {
	DefaultImpl
	Id                       string      `json:"id"`
	SessionStart             time.Time   `json:"session_start"`
	SessionEnd               interface{} `json:"session_end"`
	Duration                 int         `json:"duration"`
	MeterStart               int         `json:"meter_start"`
	MeterEnd                 int         `json:"meter_end"`
	MeterTotal               int         `json:"meter_total"`
	ChargePointId            string      `json:"chargepoint_id"`
	ConnectorUuid            string      `json:"connector_uuid"`
	ConnectorEvse            string      `json:"connector_evse"`
	CostsNet                 float64     `json:"costs_net"`
	TaxAmount                float64     `json:"tax_amount"`
	Costs                    float64     `json:"costs"`
	Currency                 string      `json:"currency"`
	Status                   string      `json:"status"`
	InitialStateOfCharge     interface{} `json:"initial_state_of_charge"`
	LastStateOfCharge        interface{} `json:"last_state_of_charge"`
	StateOfChargeLastChanged interface{} `json:"state_of_charge_last_changed"`
}

// DEFAULT

type Default interface {
	GetGAI() string
	GetProviderId() string
}

type DefaultImpl struct {
	Default
	Config *apiserver.Configuration
}

func (d *DefaultImpl) AdheresToFilter(filter [][]apiserver.FilterRule) (bool, error) {
	f := apiFilterToCommonFilter(filter)
	fp, err := utils.StructToMap(d)
	if err != nil {
		return false, fmt.Errorf("converting struct to map: %v", err)
	}
	adheres, err := common.Filter(f, fp)
	if err != nil {
		return false, err
	}
	return adheres, nil
}

func (d *DefaultImpl) GetAssetID(projectID string) (*int32, error) {
	return conf.GetAssetId(context.Background(), *d.Config, projectID, d.GetGAI())
}

func (d *DefaultImpl) SetAssetID(assetID int32, projectID string) error {
	if err := conf.InsertAsset(context.Background(), *d.Config, projectID, d.GetGAI(), assetID, d.GetProviderId()); err != nil {
		return fmt.Errorf("inserting asset to Config db: %v", err)
	}
	return nil
}

func (d *DefaultImpl) GetFunctionalChildren() []asset.FunctionalNode {
	return nil
}

func (c *Connector) GetLocationalChildren() []asset.LocationalNode {
	return nil
}

func apiFilterToCommonFilter(input [][]apiserver.FilterRule) [][]common.FilterRule {
	result := make([][]common.FilterRule, len(input))
	for i := 0; i < len(input); i++ {
		result[i] = make([]common.FilterRule, len(input[i]))
		for j := 0; j < len(input[i]); j++ {
			result[i][j] = common.FilterRule{
				Parameter: input[i][j].Parameter,
				Regex:     input[i][j].Regex,
			}
		}
	}
	return result
}
