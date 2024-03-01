package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/streamdal/streamdal/libs/protos/build/go/protos"

	"github.com/streamdal/terraform-provider-streamdal/streamdal"
)

func resourceNotification() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNotificationCreate,
		ReadContext:   resourceNotificationRead,
		UpdateContext: resourceNotificationUpdate,
		DeleteContext: resourceNotificationDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the notification configuration",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: getNotificationConfigTypes(),
			},
			"slack": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"channel": {
							Description: "The Slack channel to send the notification to",
							Type:        schema.TypeString,
							Required:    true,
						},
						"bot_token": {
							Description: "The bot token to use for sending the notification",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
			"pagerduty": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"token": {
							Description: "PagerDuty API token ",
							Type:        schema.TypeString,
							Required:    true,
						},
						"email": {
							Description: "Valid pagerduty user's email",
							Type:        schema.TypeString,
							Required:    true,
						},
						"service_id": {
							Description: "PagerDuty service's ID",
							Type:        schema.TypeString,
							Required:    true,
						},
						"urgency": {
							Description:  "The urgency of the notification",
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: getPagerDutyUrgencyTypes(),
							Default:      "low",
						},
					},
				},
			},
			"email": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description:  "Service sending the email notification",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: getEmailTypes(),
						},
						"recipients": {
							Description: "The email addresses to send the notification to",
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"from_address": {
							Description: "The email address to send the notification from",
							Type:        schema.TypeString,
							Required:    true,
						},
						"smtp": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host": {
										Description: "The SMTP server host",
										Type:        schema.TypeString,
										Required:    true,
									},
									"port": {
										Description: "The SMTP server port",
										Type:        schema.TypeInt,
										Optional:    true,
										Default:     587,
									},
									"user": {
										Description: "The SMTP server user",
										Type:        schema.TypeString,
										Required:    true,
									},
									"password": {
										Description: "The SMTP server password",
										Type:        schema.TypeString,
										Required:    true,
									},
									"use_tls": {
										Description: "Use TLS for the SMTP server",
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     true,
									},
								},
							},
						},
						"ses": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ses_region": {
										Description: "AWS region for SES service",
										Type:        schema.TypeString,
										Required:    true,
									},
									"ses_access_key": {
										Description: "AWS Access Key for SES user",
										Type:        schema.TypeString,
										Required:    true,
									},
									"ses_secret_access_key": {
										Description: "AWS Secret for SES user",
										Type:        schema.TypeString,
										Required:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceNotificationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*streamdal.Streamdal)
	resp, err := client.GetNotification(ctx, &protos.GetNotificationRequest{NotificationId: d.Id()})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.Notification.GetId())

	return diags
}

func resourceNotificationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	notification, moreDiags := buildNotification(d, m)
	if moreDiags.HasError() {
		return moreDiags
	}

	req := &protos.CreateNotificationRequest{
		Notification: notification,
	}

	resp, err := m.(*streamdal.Streamdal).CreateNotification(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.Notification.GetId())

	return diags
}

func resourceNotificationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	notification, moreDiags := buildNotification(d, m)
	if moreDiags.HasError() {
		return moreDiags
	}

	req := &protos.UpdateNotificationRequest{
		Notification: notification,
	}

	resp, err := m.(*streamdal.Streamdal).UpdateNotification(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	// TODO: need to update anything here? probably
	_ = resp

	return diags

}

func resourceNotificationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	_, err := m.(*streamdal.Streamdal).DeleteNotification(ctx, &protos.DeleteNotificationRequest{
		NotificationId: d.Id(),
	})

	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func buildNotification(d *schema.ResourceData, m interface{}) (*protos.NotificationConfig, diag.Diagnostics) {
	var diags diag.Diagnostics

	t, err := notificationConfigTypeFromString(d.Get("type").(string))
	if err != nil {
		return nil, diag.Errorf("Error converting notification type: %s", err.Error())
	}

	n := &protos.NotificationConfig{
		Name: d.Get("name").(string),
		Type: t,
	}

	switch t {
	case protos.NotificationType_NOTIFICATION_TYPE_SLACK:
		tmpSlack := d.Get("slack").([]interface{})
		if tmpSlack == nil || len(tmpSlack) == 0 {
			return nil, diag.Errorf("Error creating notification: 'slack' configuration is required")
		}

		cfg := tmpSlack[0].(map[string]interface{})
		n.Config = &protos.NotificationConfig_Slack{
			Slack: &protos.NotificationSlack{
				Channel:  cfg["channel"].(string),
				BotToken: cfg["bot_token"].(string),
			},
		}
	case protos.NotificationType_NOTIFICATION_TYPE_PAGERDUTY:
		tmpPagerDuty := d.Get("pagerduty").([]interface{})
		if tmpPagerDuty == nil || len(tmpPagerDuty) == 0 {
			return nil, diag.Errorf("Error creating notification: 'pagerduty' configuration is required")
		}

		cfg := tmpPagerDuty[0].(map[string]interface{})
		urgency, err := pagerDutyUrgencyTypeFromString(cfg["urgency"].(string))
		if err != nil {
			return nil, diag.Errorf("Error creating notification: %s", err.Error())
		}

		n.Config = &protos.NotificationConfig_Pagerduty{
			Pagerduty: &protos.NotificationPagerDuty{
				Token:     cfg["token"].(string),
				Email:     cfg["email"].(string),
				ServiceId: cfg["service_id"].(string),
				Urgency:   urgency,
			},
		}
	case protos.NotificationType_NOTIFICATION_TYPE_EMAIL:
		tmpEmail := d.Get("email").([]interface{})
		if tmpEmail == nil || len(tmpEmail) == 0 {
			return nil, diag.Errorf("Error creating notification: 'email' configuration is required")
		}

		cfg := tmpEmail[0].(map[string]interface{})
		emailType, err := emailTypeFromString(cfg["type"].(string))
		if err != nil {
			return nil, diag.Errorf("Error creating notification: %s", err.Error())
		}

		n.Config = &protos.NotificationConfig_Email{
			Email: &protos.NotificationEmail{
				Type:        emailType,
				Recipients:  interfaceToStrings(cfg["recipients"].([]interface{})),
				FromAddress: cfg["from_address"].(string),
				// Config is filled out below
			},
		}

		switch emailType {
		case protos.NotificationEmail_TYPE_SMTP:
			tmpSmtp := cfg["smtp"].([]interface{})
			if tmpSmtp == nil || len(tmpSmtp) == 0 {
				return nil, diag.Errorf("Error creating notification: 'smtp' configuration is required")
			}

			smtpCfg := tmpSmtp[0].(map[string]interface{})
			n.GetEmail().Config = &protos.NotificationEmail_Smtp{
				Smtp: &protos.NotificationEmailSMTP{
					Host:     smtpCfg["host"].(string),
					Port:     int32(smtpCfg["port"].(int)),
					User:     smtpCfg["user"].(string),
					Password: smtpCfg["password"].(string),
					UseTls:   smtpCfg["use_tls"].(bool),
				},
			}
		case protos.NotificationEmail_TYPE_SES:
			tmpSes := cfg["ses"].([]interface{})
			if tmpSes == nil || len(tmpSes) == 0 {
				return nil, diag.Errorf("Error creating notification: 'ses' configuration is required")
			}

			sesCfg := tmpSes[0].(map[string]interface{})
			n.GetEmail().Config = &protos.NotificationEmail_Ses{
				Ses: &protos.NotificationEmailSES{
					SesRegion:          sesCfg["ses_region"].(string),
					SesAccessKeyId:     sesCfg["ses_access_key"].(string),
					SesSecretAccessKey: sesCfg["ses_secret_access_key"].(string),
				},
			}

		}
	}

	return n, diags
}
