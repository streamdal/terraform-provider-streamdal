package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/streamdal/streamdal/libs/protos/build/go/protos"
	"github.com/streamdal/terraform-provider-streamdal/streamdal"
	"github.com/streamdal/terraform-provider-streamdal/util"
)

func dataSourceAudience() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAudienceRead,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"service_name": {
				Description: "The name of the service",
				Type:        schema.TypeString,
				Required:    true,
			},
			"component_name": {
				Description: "The name of the component",
				Type:        schema.TypeString,
				Required:    true,
			},
			"operation_name": {
				Description: "The name of the operation",
				Type:        schema.TypeString,
				Required:    true,
			},
			"operation_type": {
				Description: "The type of the operation, either `consumer` or `producer`",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceAudienceRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	audCfg, moreDiags := s.GetAudienceFilter(filters)
	if moreDiags.HasError() {
		return append(diags, moreDiags...)
	}

	aud := &protos.Audience{
		ServiceName:   audCfg["service_name"].(string),
		ComponentName: audCfg["component_name"].(string),
		OperationType: protos.OperationType(audCfg["operation_type"].(int)),
		OperationName: audCfg["operation_name"].(string),
	}

	d.SetId(util.AudienceToStr(aud))
	_ = d.Set("service_name", aud.ServiceName)
	_ = d.Set("component_name", aud.ComponentName)
	_ = d.Set("operation_type", aud.OperationType)
	_ = d.Set("operation_name", aud.OperationName)

	return diags
}
