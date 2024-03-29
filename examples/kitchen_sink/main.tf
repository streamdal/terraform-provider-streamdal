# This example shows how to create a multi-step pipeline and assign it to an audience.
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

# Create a notification config which notifies a slack channel
# Note: It's obviously not good practice to include bot tokens in your terraform files,
#       this is just for demonstration purposes.
resource "streamdal_notification" "slack_engineering" {
  name = "Notify Slack Engineering"
  type = "slack"
  slack {
    channel   = "engineering"
    bot_token = "xoxb-abc1234"
  }
}

# Create a pipeline with two steps
# The first step detects email fields and the second step masks the email field
# If no email addresses are detected in the payload, further processing is aborted
# If an error occurs, the pipeline is aborted and a notification is sent which includes the value of the payload
resource "streamdal_pipeline" "mask_email" {
  name = "Mask Email"

  step {
    name = "Detect Email Field"
    on_false {
      abort = "abort_current" # No need to continue this pipeline if no email found
    }
    on_error {
      abort = "abort_current"
      notification {
        notification_config_ids = [resource.streamdal_notification.slack_engineering.id]
        paths                   = []
        payload_type            = "full_payload"
      }
    }
    dynamic = false
    detective {
      type   = "pii_email"
      args   = [] # no args for this type
      negate = false
      path   = "" # No path, we will scan the entire payload
    }
  }

  step {
    name    = "Mask Email Step"
    dynamic = true
    transform {
      mask_value {
        # No path needed since dynamic=true
        # We will use the results from the first detective step
        path = ""

        # Mask the email field(s) we find with asterisks
        mask = "*"
      }
    }
  }
}

# Create audience and assign the created pipeline to it
resource "streamdal_audience" "billing_sales_report" {
  service_name   = "billing-svc"
  component_name = "kafka"
  operation_name = "gen-sales-report"
  operation_type = "consumer"
  pipeline_ids   = [resource.streamdal_pipeline.mask_email.id]
}