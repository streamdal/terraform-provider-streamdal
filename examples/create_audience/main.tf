# This example shows how to pre-create an audience in Streamdal
# Normally audiences are announced by the SDK or shim when .Process() is called
# However they can be created ahead of time in order to support creating and
# assigning a pipeline to an audience before code is ran.
terraform {
  required_providers {
    streamdal = {
      version = "0.1.2"
      source  = "streamdal/streamdal"
    }
  }
}

provider "streamdal" {
  token              = "1234"
  address            = "localhost:8082"
  connection_timeout = 10
}

resource "streamdal_audience" "billing_read_orders" {
  service_name   = "billing-svc"
  component_name = "kafka"
  operation_name = "read_orders"
  operation_type = "consumer"
}