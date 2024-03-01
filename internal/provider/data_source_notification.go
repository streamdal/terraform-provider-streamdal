package provider

import (
	"context"

	"github.com/streamdal/terraform-provider-streamdal/streamdal"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNotification() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceNotificationRead,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"id": {
				Description: "Notification Config ID",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Pipeline name",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceNotificationRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var filters []*streamdal.Filter

	s := m.(*streamdal.Streamdal)

	if v, ok := d.GetOk("filter"); ok {
		filters = buildFiltersDataSource(v.(*schema.Set))
	} else {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "No filters defined",
			Detail:   "At least one filter must be defined",
		})
	}

	notificationCfg, moreDiags := s.GetNotificationConfigFilter(filters)
	if moreDiags.HasError() {
		return append(diags, moreDiags...)
	}

	d.SetId(notificationCfg["id"].(string))
	_ = d.Set("name", notificationCfg["name"].(string))

	return diags
}
