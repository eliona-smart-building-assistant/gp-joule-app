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
	"math"
	"time"

	"github.com/eliona-smart-building-assistant/go-eliona/asset"
	"github.com/eliona-smart-building-assistant/go-eliona/utils"
	"github.com/eliona-smart-building-assistant/go-utils/common"
)

// ROOT

type Root struct {
	Clusters []*Cluster

	// own attributes
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
	if err := conf.InsertAsset(context.Background(), r.Config, projectID, r.GetGAI(), assetID, r.GetAssetType(), "", ""); err != nil {
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

	// own attributes
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
	if err := conf.InsertAsset(context.Background(), c.Config, projectID, c.GetGAI(), assetID, c.GetAssetType(), "", c.Name); err != nil {
		return fmt.Errorf("inserting asset to Config db: %v", err)
	}
	return nil
}

func (c *Cluster) GetLocationalChildren() []asset.LocationalNode {
	locationalChildren := make([]asset.LocationalNode, 0)
	for _, chargingPoint := range c.ChargePoints {
		chargingPoint.Cluster = c
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
	ConnectorsOccupied  int          `json:"connectors_occupied" eliona:"connectors_occupied" subtype:"status"`
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

	// own attributes
	Cluster *Cluster
	Config  *apiserver.Configuration
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
	if err := conf.InsertAsset(context.Background(), cp.Config, projectID, cp.GetGAI(), assetID, cp.GetAssetType(), cp.Cluster.GetName(), cp.ChargePointId); err != nil {
		return fmt.Errorf("inserting asset to Config db: %v", err)
	}
	return nil
}

func (cp *ChargePoint) GetLocationalChildren() []asset.LocationalNode {
	locationalChildren := make([]asset.LocationalNode, 0)

	// Add connectors
	for idx, connector := range cp.Connectors {
		connector.ChargePoint = cp
		connector.Config = cp.Config
		connector.Index = idx + 1
		if connector.ChargingSession != nil && connector.ChargingSession.MeterTotal > 0 && connector.ChargingSession.SessionStart != nil && connector.ChargingSession.SessionEnd != nil {
			connector.Duration = connector.ChargingSession.Duration
			connector.MeterTotal = int(math.Max(float64(connector.ChargingSession.MeterTotal), 0))
			connector.Occupied = mapOccupancyStatus(connector.Status)
		} else {
			connector.Occupied = mapOccupancyStatus("available")
		}
		locationalChildren = append(locationalChildren, connector)
	}

	return locationalChildren
}

// CONNECTOR

type Connector struct {
	ConnectorId     string           `json:"uuid"`
	EvseId          string           `json:"evseid"`
	Status          string           `json:"status" eliona:"status" subtype:"status"`
	MaxPower        int              `json:"max_power" eliona:"max_power" subtype:"info"`
	ChargePointType string           `json:"chargepoint_type"`
	PlugType        string           `json:"plug_type" eliona:"plug_type,filterable"`
	ChargingSession *ChargingSession `json:"charging_session"`

	// own attributes
	ChargePoint *ChargePoint
	Config      *apiserver.Configuration
	MeterTotal  int `eliona:"current_energy" subtype:"input"`
	Duration    int `eliona:"current_duration" subtype:"input"`
	Occupied    int `eliona:"occupied" subtype:"status"`
	Index       int
}

func (c *Connector) GetName() string {
	return fmt.Sprintf("%s %s %d", c.PlugType, c.ChargePointType, c.Index)
}

func (c *Connector) GetDescription() string {
	return c.EvseId
}

func (c *Connector) GetAssetType() string {
	return "gp_joule_connector"
}

func (c *Connector) GetGAI() string {
	return c.GetAssetType() + "_" + c.ConnectorId
}

func (c *Connector) AdheresToFilter(filter [][]apiserver.FilterRule) (bool, error) {
	return adheresToFilter(c, filter)
}

func (c *Connector) GetAssetID(projectID string) (*int32, error) {
	return conf.GetAssetId(context.Background(), c.Config, projectID, c.GetGAI())
}

func (c *Connector) SetAssetID(assetID int32, projectID string) error {
	if err := conf.InsertAsset(context.Background(), c.Config, projectID, c.GetGAI(), assetID, c.GetAssetType(), c.ChargePoint.ChargePointId, c.ConnectorId); err != nil {
		return fmt.Errorf("inserting asset to Config db: %v", err)
	}
	return nil
}

func (c *Connector) GetLocationalChildren() []asset.LocationalNode {
	locationalChildren := make([]asset.LocationalNode, 0)

	// Add one sessions container
	locationalChildren = append(locationalChildren, &SessionsLog{
		Connector: c,
		Config:    c.Config,
	})

	return locationalChildren
}

// COMPLETED SESSIONS

type SessionsLog struct {
	Connector *Connector
	Config    *apiserver.Configuration
}

func (cs *SessionsLog) GetName() string {
	return fmt.Sprintf("%s session log", cs.Connector.GetName())
}

func (cs *SessionsLog) GetDescription() string {
	return fmt.Sprintf("Session log for %s", cs.Connector.GetName())
}

func (cs *SessionsLog) GetAssetType() string {
	return "gp_joule_session_log"
}

func (cs *SessionsLog) GetGAI() string {
	return cs.GetAssetType() + "_" + cs.Connector.ConnectorId
}

func (cs *SessionsLog) AdheresToFilter(filter [][]apiserver.FilterRule) (bool, error) {
	return adheresToFilter(cs, filter)
}

func (cs *SessionsLog) GetAssetID(projectID string) (*int32, error) {
	return conf.GetAssetId(context.Background(), cs.Config, projectID, cs.GetGAI())
}

func (cs *SessionsLog) SetAssetID(assetID int32, projectID string) error {
	if err := conf.InsertAsset(context.Background(), cs.Config, projectID, cs.GetGAI(), assetID, cs.GetAssetType(), cs.Connector.ConnectorId, ""); err != nil {
		return fmt.Errorf("inserting asset to Config db: %v", err)
	}
	return nil
}

func (cs *SessionsLog) GetLocationalChildren() []asset.LocationalNode {
	return make([]asset.LocationalNode, 0)
}

// SESSION

type ChargingSession struct {
	Id                       string      `json:"id"`
	SessionStart             *time.Time  `json:"session_start"`
	SessionEnd               *time.Time  `json:"session_end"`
	Duration                 int         `json:"duration"`
	MeterStart               int         `json:"meter_start"`
	MeterEnd                 int         `json:"meter_end"`
	MeterTotal               int         `json:"meter_total"`
	ChargePointId            string      `json:"chargepoint_id"`
	ConnectorId              string      `json:"connector_uuid"`
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

// ERROR

type ErrorNotification struct {
	Id            string     `json:"id"`
	ChargePointId string     `json:"chargepoint_id"`
	ConnectorId   *string    `json:"connector_uuid"`
	ErrorCode     string     `json:"error_code"`
	ErrorInfo     string     `json:"error_info"`
	VendorCode    string     `json:"vendor_code"`
	OccurredAt    *time.Time `json:"occurred_at"`
	ResolvedAt    *time.Time `json:"resolved_at"`
}

// HELPER

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

func mapOccupancyStatus(status string) int {
	if status == "available" {
		return 0
	}
	if status == "occupied" {
		return 1
	}
	return -1
}
