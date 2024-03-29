---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "streamdal Provider"
subcategory: ""
description: |-
  
---

# Streamdal Provider

![logo](https://github.com/streamdal/streamdal/raw/main/assets/img/streamdal-logo-light.png#gh-light-mode-only)

[![Discord](https://img.shields.io/badge/Community-Discord-4c57e8.svg)](https://discord.gg/streamdal)

This provider is used to interact with the [Streamdal server API](https://github.com/streamdal/streamdal)


<!-- schema generated by tfplugindocs -->
## Provider Schema

| variable | type   | description                                   | envar |
|:---|:-------|:----------------------------------------------|:---|
| address | string | The address of your Streamdal server install. | `STREAMDAL_ADDRESS` |
| connection_timeout | int    | gRPC connection attempt timeout in seconds.   | `STREAMDAL_CONNECTION_TIMEOUT` |
| token | string | API Auth Token                                | `STREAMDAL_TOKEN` |

## Example Provider Setup

```hcl
terraform {
  required_providers {
    streamdal = {
      source = "streamdal/streamdal"
      version = "0.1.1"
    }
  }
}

provider "streamdal" {
  token              = "1234"
  address            = "localhost:8082"
  connection_timeout = 10
}
```