---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "streamdal_notification Resource - terraform-provider-streamdal"
subcategory: ""
description: |-
  
---

# streamdal_notification (Resource)

The `streamdal_notification` resource allows you to create, assign, and delete notification configurations that
can be used inside on_true, on_false, and on_error blocks of a pipeline step

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

resource "streamdal_notification" "slack_engineering" {
  name = "Notify Slack Engineering"
  type = "slack"
  slack {
    channel   = "engineering"
    bot_token = "1234"
  }
}
```


<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **name** (String)
- **type** (String)

One of the following notification type blocks must be set:

- **email** (Block) (see [below for nested schema](#nestedblock--email))
- **pagerduty** (Block) (see [below for nested schema](#nestedblock--pagerduty))
- **slack** (Block) (see [below for nested schema](#nestedblock--slack))

### Read-Only

- **id** (String) The ID of the notification configuration

<a id="nestedblock--email"></a>
### Nested Schema for `email`

Required:

- **from_address** (String) The email address to send the notification from
- **recipients** (List of String) The email addresses to send the notification to
- **type** (String) Service sending the email notification

One email type block must be specified (either `ses` or `smtp`):

- **ses** (Block) (see [below for nested schema](#nestedblock--email--ses))
- **smtp** (Block) (see [below for nested schema](#nestedblock--email--smtp))

<a id="nestedblock--email--ses"></a>
### Nested Schema for `email.ses`

Required:

- **ses_access_key** (String) AWS Access Key for SES user
- **ses_region** (String) AWS region for SES service
- **ses_secret_access_key** (String) AWS Secret for SES user


<a id="nestedblock--email--smtp"></a>
### Nested Schema for `email.smtp`

Required:

- **host** (String) The SMTP server host
- **password** (String) The SMTP server password
- **user** (String) The SMTP server user

Optional:

- **port** (Number) The SMTP server port (Default: `587`)
- **use_tls** (Boolean) Use TLS for the SMTP server (Default: `true`)



<a id="nestedblock--pagerduty"></a>
### Nested Schema for `pagerduty`

Required:

- **email** (String) Valid pagerduty user's email
- **service_id** (String) PagerDuty service's ID
- **token** (String) PagerDuty API token

Optional:

- **urgency** (String) The urgency of the notification (Default: `low`)


<a id="nestedblock--slack"></a>
### Nested Schema for `slack`

Required:

- **bot_token** (String) The bot token to use for sending the notification
- **channel** (String) The Slack channel to send the notification to


