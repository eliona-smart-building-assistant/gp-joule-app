package main

import (
	"github.com/eliona-smart-building-assistant/app-integration-tests/app"
	"github.com/eliona-smart-building-assistant/app-integration-tests/assert"
	"github.com/eliona-smart-building-assistant/app-integration-tests/test"
	"testing"
)

func TestApp(t *testing.T) {
	app.StartApp()
	test.AppWorks(t)
	t.Run("TestSchema", schema)
	t.Run("TestAssetTypes", assetTypes)
	t.Run("TestWidgetTypes", widgetTypes)
	app.StopApp()
}

func assetTypes(t *testing.T) {
	t.Parallel()

	assert.AssetTypeExists(t, "gp_joule_charge_point", []string{"model"})
	assert.AssetTypeExists(t, "gp_joule_cluster", []string{})
	assert.AssetTypeExists(t, "gp_joule_connector", []string{"status"})
	assert.AssetTypeExists(t, "gp_joule_root", []string{})
	assert.AssetTypeExists(t, "gp_joule_session_log", []string{"energy"})
}

func widgetTypes(t *testing.T) {
	t.Parallel()

	assert.WidgetTypeExists(t, "gp_joule_connector")
}

func schema(t *testing.T) {
	t.Parallel()

	assert.SchemaExists(t, "gp_joule", []string{"configuration", "asset"})
}
