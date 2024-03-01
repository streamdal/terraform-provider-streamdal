# Testing locally

1. `make`
2. Run TF files with the following provider config:
   ```hcl
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
   ```
   
TODO: expand