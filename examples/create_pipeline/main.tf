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

  # Change to your schema name or use default JSON
  # Wildcards "*" are accepted
  default = "test slack cfg"
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

#resource "streamdal_pipeline" "detect_hostname" {
#  name = "Step 1"
#  step {
#    name = "Delete Field"
#    on_false {
#      abort="abort_current"
#    }
#    on_error {
#      abort="abort_all"
#      notification {
#        notification_config_ids = ["958a663a-561f-4463-acc7-d84ab2043c09"]
#        paths = ["object.payload"]
#        payload_type = "select_paths"
#      }
#
#    }
#    dynamic=false
#    detective {
#      type = "hostname"
#      args = [] # no args for this type
#      negate = false
#      path = "object.payload"
#    }
#  }
#  step {
#    name = "Replace Field Value Step"
#    dynamic=false
#    transform {
#      type = "replace_value" # TODO: can we eliminate this?
#      replace_value {
#        path = "object.payload"
#        value = "\"omg replaced!\""
#      }
#    }
#  }
#}

# TODO: need a notification config data-source and resource
# TODO: need to pass the notification ID from the filter to the detect_hostname step