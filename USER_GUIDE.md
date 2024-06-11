
# GP Joule User Guide

### Introduction

> The GP Joule app provides integration and synchronization between Eliona and GP Joule services.

## Overview

This guide provides instructions on configuring, installing, and using the GP Joule app to manage resources and synchronize data between Eliona and GP Joule services.

## Installation

Install the GP Joule app via the Eliona App Store.

## Configuration

The GP Joule app requires configuration through Eliona’s settings interface. Below are the general steps and details needed to configure the app effectively.

### Registering the app in GP Joule Service

Create credentials in GP Joule Service to connect the GP Joule services from Eliona. All required credentials are listed below in the [configuration section](#configure-the-gp-joule-app).  

Please contact GP Joule to get the URL for API access and your access key.  

### Configure the GP Joule app 

Configurations can be created in Eliona under `Apps > GP Joule > Settings` which opens the app's [Generic Frontend](https://doc.eliona.io/collection/v/eliona-english/manuals/settings/apps). Here you can use the appropriate endpoint with the POST method. Each configuration requires the following data:

| Attribute         | Description                                                                     |
|-------------------|---------------------------------------------------------------------------------|
| `rootUrl`         | URL of the GP Joule API services.                                               |
| `apiKey`          | Client secrets obtained from the GP Joule service.                              |
| `assetFilter`     | Filtering asset during [Continuous Asset Creation](#continuous-asset-creation). |
| `enable`          | Flag to enable or disable this configuration.                                   |
| `refreshInterval` | Interval in seconds for data synchronization.                                   |
| `requestTimeout`  | API query timeout in seconds.                                                   |
| `projectIDs`      | List of Eliona project IDs for data collection.                                 |

Example configuration JSON:

```json
{
  "rootUrl": "http://gpjoule.api.url/",
  "apiKey": "s3cr3t",
  "enable": true,
  "refreshInterval": 60,
  "requestTimeout": 120,
  "projectIDs": [
    "10"
  ]
}
```

## Continuous Asset Creation

Once configured, the app starts Continuous Asset Creation (CAC). Discovered resources are automatically created as assets in Eliona, and users are notified via Eliona’s notification system.

The GP Joule infrastructure is managed through the charge point asset, which groups all charge points by their cluster name. Each charge point can have one or more connectors, each with its own charging properties.

Additionally, a session log asset is created for each connector, providing historical records of all charging sessions.

## Additional Features

### Dashboard templates

The app offers a predefined dashboard that clearly displays the most important information. YOu can create such a dashboard under `Dashboards > Copy Dashboard > From App > GP Joule`.
