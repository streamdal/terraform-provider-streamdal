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

variable "notification_config_id" {
  type = string

  # Wildcards "*" are accepted
  default = "test slack *"
}

data "streamdal_notification" "slack_test" {
  filter {
    name   = "name"
    values = [var.notification_config_id]
  }
}

output "notification" {
  value = data.streamdal_notification.slack_test
}
