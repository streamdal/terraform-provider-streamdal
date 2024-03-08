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

resource "streamdal_notification" "slack_engineering" {
  name = "Notify Slack Engineering"
  type = "slack"
  slack {
    channel   = "engineering"
    bot_token = "1234"
  }
}