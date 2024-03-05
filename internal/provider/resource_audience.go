package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/streamdal/streamdal/libs/protos/build/go/protos"
	"github.com/streamdal/terraform-provider-streamdal/streamdal"
	"github.com/streamdal/terraform-provider-streamdal/util"
)

func resourceAudience() *schema.Resource {
	sch := audienceSchema()
	sch["pipeline_ids"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}

	return &schema.Resource{
		CreateContext: resourceAudienceCreate,
		ReadContext:   resourceAudienceRead,
		UpdateContext: resourceAudienceUpdate,
		DeleteContext: resourceAudienceDelete,

		Schema: sch,
	}
}

func resourceAudienceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*streamdal.Streamdal)

	aud := &protos.Audience{
		ServiceName:   d.Get("service_name").(string),
		ComponentName: d.Get("component_name").(string),
		OperationType: audienceOperationTypeFromString(d.Get("operation_type").(string)),
		OperationName: d.Get("operation_name").(string),
	}

	if _, err := client.CreateAudience(ctx, &protos.CreateAudienceRequest{Audience: aud}); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(util.AudienceToStr(aud))

	// Assign pipelines
	pipelineIDs := interfaceToStrings(d.Get("pipeline_ids"))
	if _, err := client.SetPipelines(ctx, aud, pipelineIDs); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

// TODO: both GetAudience() and GetPipelinesForAudience() call GetAll(). Let's see if we can combine them
func resourceAudienceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*streamdal.Streamdal)

	aud, err := client.GetAudience(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Get pipeline assignments. These come from a GetAll() call
	pipelineIDs, err := client.GetPipelinesForAudience(ctx, aud)
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("service_name", aud.ServiceName)
	_ = d.Set("component_name", aud.ComponentName)
	_ = d.Set("operation_name", aud.OperationName)
	_ = d.Set("operation_type", audienceOperationTypeToString(aud.OperationType))
	_ = d.Set("pipeline_ids", pipelineIDs)

	return diags
}

func resourceAudienceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*streamdal.Streamdal)

	// Verify audience exists, otherwise error out.
	// Audiences only support updates to pipeline assignments, not the actual audience data itself.
	aud, err := client.GetAudience(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Update pipeline assignments
	pipelineIDs := interfaceToStrings(d.Get("pipeline_ids"))

	if _, err := client.SetPipelines(ctx, aud, pipelineIDs); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceAudienceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	aud := util.AudienceFromStr(d.Id())

	client := m.(*streamdal.Streamdal)
	if _, err := client.DeleteAudience(ctx, &protos.DeleteAudienceRequest{Audience: aud}); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
