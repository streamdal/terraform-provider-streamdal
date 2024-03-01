package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/streamdal/terraform-provider-streamdal/streamdal"
)

func dataSourcePipeline() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourcePipelineRead,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"id": {
				Description: "Pipeline ID",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Pipeline Name",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourcePipelineRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var filters []*streamdal.Filter

	client := m.(*streamdal.Streamdal)

	if v, ok := d.GetOk("filter"); ok {
		filters = buildFiltersDataSource(v.(*schema.Set))
	} else {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "No filters defined",
			Detail:   "At least one filter must be defined",
		})
	}

	pipeline, moreDiags := client.GetPipelineFilter(filters)
	if moreDiags.HasError() {
		return append(diags, moreDiags...)
	}

	d.SetId(pipeline["id"].(string))
	_ = d.Set("name", pipeline["name"].(string))

	return diags
}
