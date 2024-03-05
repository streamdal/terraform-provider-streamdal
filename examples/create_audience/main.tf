# This example shows how to pre-create an audience in Streamdal
# Normally audiences are announced by the SDK or shim when .Process() is called
# However they can be created ahead of time in order to support creating and
# assigning a pipeline to an audience before code is ran.
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

resource "streamdal_audience" "testaud" {
  service_name   = "test_service"
  component_name = "kafka"
  operation_name = "read_stuff3"
  operation_type = "consumer"
}