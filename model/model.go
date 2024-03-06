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
	"strings"
	"time"

	"github.com/eliona-smart-building-assistant/go-eliona/asset"
	"github.com/eliona-smart-building-assistant/go-eliona/utils"
	"github.com/eliona-smart-building-assistant/go-utils/common"
)

// ROOT

type Root struct {
	Clusters []*Cluster

	Config *apiserver.Configuration
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

func (r *Root) AdheresToFilter(filter [][]apiserver.FilterRule) (bool, error) {
	return adheresToFilter(r, filter)
}

func (r *Root) GetAssetID(projectID string) (*int32, error) {
	return conf.GetAssetId(context.Background(), r.Config, projectID, r.GetGAI())
}

func (r *Root) SetAssetID(assetID int32, projectID string) error {
	if err := conf.InsertAsset(context.Background(), r.Config, projectID, r.GetGAI(), assetID, r.GetAssetType(), ""); err != nil {
		return fmt.Errorf("inserting asset to Config db: %v", err)
	}
	return nil
}

func (r *Root) GetFunctionalChildren() []asset.FunctionalNode {
	return make([]asset.FunctionalNode, 0)
}

func (r *Root) GetLocationalChildren() []asset.LocationalNode {
	locationalChildren := make([]asset.LocationalNode, 0)
	for _, cluster := range r.Clusters {
		cluster.Config = r.Config
		locationalChildren = append(locationalChildren, cluster)
	}
	return locationalChildren
}

// CLUSTER

type Cluster struct {
	Name         string         `json:"name" eliona:"name,filterable"`
	ChargePoints []*ChargePoint `json:"chargepoints"`

	Config *apiserver.Configuration
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

func (c *Cluster) AdheresToFilter(filter [][]apiserver.FilterRule) (bool, error) {
	return adheresToFilter(c, filter)
}

func (c *Cluster) GetAssetID(projectID string) (*int32, error) {
	return conf.GetAssetId(context.Background(), c.Config, projectID, c.GetGAI())
}

func (c *Cluster) SetAssetID(assetID int32, projectID string) error {
	if err := conf.InsertAsset(context.Background(), c.Config, projectID, c.GetGAI(), assetID, c.GetAssetType(), c.Name); err != nil {
		return fmt.Errorf("inserting asset to Config db: %v", err)
	}
	return nil
}

func (c *Cluster) GetLocationalChildren() []asset.LocationalNode {
	locationalChildren := make([]asset.LocationalNode, 0)
	for _, chargingPoint := range c.ChargePoints {
		chargingPoint.Config = c.Config
		locationalChildren = append(locationalChildren, chargingPoint)
	}
	return locationalChildren
}

// CHARGE POINT

type ChargePoint struct {
	ChargePointId       string       `json:"chargepoint_id" eliona:"id,filterable"`
	ChargePointOcppId   string       `json:"chargepoint_ocpp_id"`
	Name                string       `json:"name" eliona:"name,filterable"`
	NameInternal        string       `json:"name_internal" eliona:"name_internal,filterable"`
	Status              string       `json:"status" eliona:"status" subtype:"status"`
	CommunicationStatus int          `json:"communication_status"`
	ConnectorsTotal     int          `json:"connectors_total" eliona:"connectors_total" subtype:"info"`
	ConnectorsFree      int          `json:"connectors_free"`
	ConnectorsFaulted   int          `json:"connectors_faulted"`
	ConnectorsOccupied  int          `json:"connectors_occupied" eliona:"connectors_occupied" subtype:"input"`
	Manufacturer        string       `json:"manufacturer" eliona:"manufacturer,filterable" subtype:"info"`
	Model               string       `json:"model" eliona:"model,filterable" subtype:"info"`
	Lat                 float64      `json:"lat"`
	Long                float64      `json:"long"`
	Street              string       `json:"street"`
	Zip                 int          `json:"zip"`
	City                string       `json:"city"`
	CountryCode         interface{}  `json:"country_code"`
	Country             string       `json:"country"`
	Error               interface{}  `json:"error"`
	Connectors          []*Connector `json:"connectors"`

	Config *apiserver.Configuration
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

func (cp *ChargePoint) AdheresToFilter(filter [][]apiserver.FilterRule) (bool, error) {
	return adheresToFilter(cp, filter)
}

func (cp *ChargePoint) GetAssetID(projectID string) (*int32, error) {
	return conf.GetAssetId(context.Background(), cp.Config, projectID, cp.GetGAI())
}

func (cp *ChargePoint) SetAssetID(assetID int32, projectID string) error {
	if err := conf.InsertAsset(context.Background(), cp.Config, projectID, cp.GetGAI(), assetID, cp.GetAssetType(), cp.ChargePointId); err != nil {
		return fmt.Errorf("inserting asset to Config db: %v", err)
	}
	return nil
}

func (cp *ChargePoint) GetLocationalChildren() []asset.LocationalNode {
	locationalChildren := make([]asset.LocationalNode, 0)

	// Add connectors
	for _, connector := range cp.Connectors {
		connector.Config = cp.Config
		locationalChildren = append(locationalChildren, connector)
	}

	// Add one sessions container
	locationalChildren = append(locationalChildren, &ChargingSessions{
		ChargePoint: cp,
		Config:      cp.Config,
	})

	return locationalChildren
}

// SESSIONS

type ChargingSessions struct {
	ChargePoint *ChargePoint
	Config      *apiserver.Configuration
}

func (s *ChargingSessions) GetName() string {
	return fmt.Sprintf("%s sessions", s.ChargePoint.Name)
}

func (s *ChargingSessions) GetDescription() string {
	return fmt.Sprintf("ChargingSessions for %s", s.ChargePoint.NameInternal)
}

func (s *ChargingSessions) GetAssetType() string {
	return "gp_joule_resent_sessions"
}

func (s *ChargingSessions) GetGAI() string {
	return s.GetAssetType() + "_" + s.ChargePoint.ChargePointId
}

func (s *ChargingSessions) AdheresToFilter(filter [][]apiserver.FilterRule) (bool, error) {
	return adheresToFilter(s, filter)
}

func (s *ChargingSessions) GetAssetID(projectID string) (*int32, error) {
	return conf.GetAssetId(context.Background(), s.Config, projectID, s.GetGAI())
}

func (s *ChargingSessions) SetAssetID(assetID int32, projectID string) error {
	if err := conf.InsertAsset(context.Background(), s.Config, projectID, s.GetGAI(), assetID, s.GetAssetType(), s.ChargePoint.ChargePointId); err != nil {
		return fmt.Errorf("inserting asset to Config db: %v", err)
	}
	return nil
}

func (s *ChargingSessions) GetLocationalChildren() []asset.LocationalNode {
	return make([]asset.LocationalNode, 0)
}

// CONNECTOR

type Connector struct {
	Uuid            string `json:"uuid"`
	EvseId          string `json:"evseid"`
	Status          string `json:"status" eliona:"status" subtype:"status"`
	MaxPower        int    `json:"max_power" eliona:"max_power" subtype:"info"`
	ChargePointType string `json:"chargepoint_type"`
	PlugType        string `json:"plug_type" eliona:"plug_type,filterable"`

	Config *apiserver.Configuration
}

func (c *Connector) GetName() string {
	return strings.Trim(c.PlugType+" "+c.ChargePointType, " ")
}

func (c *Connector) GetDescription() string {
	return c.ChargePointType
}

func (c *Connector) GetAssetType() string {
	return "gp_joule_connector"
}

func (c *Connector) GetGAI() string {
	return c.GetAssetType() + "_" + c.Uuid
}

func (c *Connector) AdheresToFilter(filter [][]apiserver.FilterRule) (bool, error) {
	return adheresToFilter(c, filter)
}

func (c *Connector) GetAssetID(projectID string) (*int32, error) {
	return conf.GetAssetId(context.Background(), c.Config, projectID, c.GetGAI())
}

func (c *Connector) SetAssetID(assetID int32, projectID string) error {
	if err := conf.InsertAsset(context.Background(), c.Config, projectID, c.GetGAI(), assetID, c.GetAssetType(), c.Uuid); err != nil {
		return fmt.Errorf("inserting asset to Config db: %v", err)
	}
	return nil
}

func (c *Connector) GetLocationalChildren() []asset.LocationalNode {
	return make([]asset.LocationalNode, 0)
}

// SESSION

type ChargingSession struct {
	Id                       string      `json:"id"`
	SessionStart             time.Time   `json:"session_start"`
	SessionEnd               time.Time   `json:"session_end"`
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

func adheresToFilter[T any](data *T, filter [][]apiserver.FilterRule) (bool, error) {
	f := apiFilterToCommonFilter(filter)
	fp, err := utils.StructToMap(data)
	if err != nil {
		return false, fmt.Errorf("converting struct to map: %v", err)
	}
	adheres, err := common.Filter(f, fp)
	if err != nil {
		return false, err
	}
	return adheres, nil
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
