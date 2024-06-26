# Eliona App for GP Joule

This app connects [GP Joule API](https://www.gp-joule.com/) and gathers information about data of charging stations.

## Configuration

The app needs environment variables and database tables for configuration. To edit the database tables the app provides an own API access.


### Registration in Eliona ###

To start and initialize an app in an Eliona environment, the app has to be registered in Eliona. For this, entries in database tables `public.eliona_app` and `public.eliona_store` are necessary.

This initialization can be handled by the `reset.sql` script.


### Environment variables

- `CONNECTION_STRING`: configures the [Eliona database](https://github.com/eliona-smart-building-assistant/go-eliona/tree/main/db). Otherwise, the app can't be initialized and started (e.g. `postgres://user:pass@localhost:5432/iot`).

- `INIT_CONNECTION_STRING`: configures the [Eliona database](https://github.com/eliona-smart-building-assistant/go-eliona/tree/main/db) for app initialization like creating schema and tables (e.g. `postgres://user:pass@localhost:5432/iot`). Default is content of `CONNECTION_STRING`.

- `API_ENDPOINT`:  configures the endpoint to access the [Eliona API v2](https://github.com/eliona-smart-building-assistant/eliona-api). Otherwise, the app can't be initialized and started. (e.g. `http://api-v2:3000/v2`)

- `API_TOKEN`: defines the secret to authenticate the app and access the Eliona API.

- `API_SERVER_PORT`(optional): define the port the API server listens. The default value is Port `3000`.

- `LOG_LEVEL`(optional): defines the minimum level that should be [logged](https://github.com/eliona-smart-building-assistant/go-utils/blob/main/log/README.md). The default level is `info`.


### Database tables ###

The app requires configuration data that remains in the database. To do this, the app creates its own database schema `gp_joule` during initialization. To modify and handle the configuration data the app provides an API access. Have a look at the [API specification](https://eliona-smart-building-assistant.github.io/open-api-docs/?https://raw.githubusercontent.com/eliona-smart-building-assistant/gp-joule-app/develop/openapi.yaml) how the configuration tables should be used.

- `gp_joule.configuration`: Contains configuration of the app. Editable through the API.

- `gp_joule.asset`: Provides asset mapping. Maps broker's asset IDs to Eliona asset IDs.

**Generation**: to generate access method to database see Generation section below.


## References

### App API ###

The app provides its own API to access configuration data and other functions. The full description of the API is defined in the `openapi.yaml` OpenAPI definition file.

- [API Reference](https://eliona-smart-building-assistant.github.io/open-api-docs/?https://raw.githubusercontent.com/eliona-smart-building-assistant/gp-joule-app/develop/openapi.yaml) shows details of the API

**Generation**: to generate api server stub see Generation section below.

### Configuring the app ###

To use the app it is necessary to create at least one configuration. A configuration points to one GP-Joule API.
Also, a configuration defines at least one Eliona project ID for automatic asset creation.

A minimum configuration that can used by the app's API endpoint `POST /configs` is:

```json
{
    "rootUrl": "https://api.joule.net/example/url",
    "apiKey": "secret",
    "enable": true,
    "projectIDs": ["10"]
}
```

### Eliona assets ###

This app creates Eliona asset types and attribute sets during initialization.

The data is written for each charging station, structured into different subtypes of Elinoa assets. The following subtypes are defined:

- `Info`: Static data which provides information about a device like address and firmware info.
- `Status`: Device status information, like battery level.
- `Input`: Current values reported by sensors.
- `Output`: Values that are to be passed back to the provider.

### Continuous asset creation ###

Assets for all devices connected to the GP Joule API are created automatically when a configuration is added.

To select which assets to create, a filter could be specified in config. The schema of the filter is defined in the `openapi.yaml` file.

Possible filter parameters are defined in the structs in `model.go` and marked with `eliona:"attribute_name,filterable"` field tag.

To avoid conflicts, the Global Asset Identifier is a manufacturer's ID prefixed with asset type name as a namespace.

### Dashboard ###

An example dashboard meant for a quick start or showcasing the apps abilities can be obtained by accessing the dashboard endpoint defined in the `openapi.yaml` file.

## Tools

### Generate API server stub ###

For the API server the [OpenAPI Generator](https://openapi-generator.tech/docs/generators/openapi-yaml) for go-server is used to generate a server stub. The easiest way to generate the server files is to use one of the predefined generation script which use the OpenAPI Generator Docker image.

```
.\generate-api-server.cmd # Windows
./generate-api-server.sh # Linux
```

### Generate Database access ###

For the database access [SQLBoiler](https://github.com/volatiletech/sqlboiler) is used. The easiest way to generate the database files is to use one of the predefined generation script which use the SQLBoiler implementation.

```
.\generate-db.cmd # Windows
./generate-db.sh # Linux
```

### Generate asset type descriptions ###

For generating asset type descriptions from field-tag-annotated structs, [asse-from-struct tool](https://github.com/eliona-smart-building-assistant/dev-utilities) can be used.
