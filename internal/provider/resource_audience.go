package provider

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/streamdal/streamdal/libs/protos/build/go/protos"
	"github.com/streamdal/terraform-provider-streamdal/streamdal"
	"github.com/streamdal/terraform-provider-streamdal/util"
)

func resourceAudience() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAudienceCreate,
		ReadContext:   resourceAudienceRead,
		UpdateContext: resourceAudienceUpdate,
		DeleteContext: resourceAudienceDelete,

		Schema: audienceSchema(),
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

	return diags
}

func resourceAudienceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*streamdal.Streamdal)

	aud, err := client.GetAudience(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("service_name", aud.ServiceName)
	_ = d.Set("component_name", aud.ComponentName)
	_ = d.Set("operation_name", aud.OperationName)
	_ = d.Set("operation_type", audienceOperationTypeToString(aud.OperationType))

	return diags

}

func resourceAudienceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// TODO: implement if possible
	return diag.FromErr(errors.New("unimplemented"))
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
