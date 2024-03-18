# This example shows how to create a step which simply validates a JSON payload
terraform {
  required_providers {
    streamdal = {
      version = "0.1.0"
      source  = "streamdal.com/tf/streamdal"
    }
  }
}

provider "streamdal" {
  token              = "1234"
  address            = "localhost:8082"
  connection_timeout = 10
}

resource "streamdal_pipeline" "validate_my_json" {
  name = "Validate JSON Payload"
  step {
    name    = "Validate JSON Step"
    dynamic = true
    valid_json {
    }
  }
}