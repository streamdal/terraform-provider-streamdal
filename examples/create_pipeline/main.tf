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

resource "streamdal_pipeline" "detect_email" {
  name = "Mask Email"

  step {
    name = "Detect Email Field"
    on_false {
      abort="abort_current"
    }
    on_error {
      abort="abort_all"
      notification {
        notification_config_ids = ["958a663a-561f-4463-acc7-d84ab2043c09"]
        paths = ["object.payload"]
        payload_type = "select_paths"
      }

    }
    dynamic=false
    detective {
      type = "pii_email"
      args = [] # no args for this type
      negate = false
      path = "object.payload"
    }
  }

  step {
    name = "Replace Field Value Step"
    dynamic=true
    transform {
      type = "mask_value" # TODO: can we eliminate this?
      mask_value {
        path = "email"
        mask = "*"
      }
    }
  }
}
