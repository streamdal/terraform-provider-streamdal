package provider

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func audienceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
			Description:  "The type of the operation, either `consumer` or `producer`",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: getAudienceOperationTypes(),
		},
	}
}

// conditionSchema returns the schema for a PipelineStepCondition message
// This is used to define the schema for the on_true, on_false, and on_error fields of a PipelineStep
func conditionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"abort": {
				Description:  "Abort",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "unset",
				ValidateFunc: getAbortConditions(),
			},
			"metadata": {
				Description: "Metadata",
				Type:        schema.TypeMap,
				Optional:    true,
				Default:     map[string]interface{}{},
			},
			"notification": {
				Description: "Notification Config",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"notification_config_ids": {
							Description: "Notification Config IDs",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							DefaultFunc: func() (interface{}, error) {
								return []string{}, nil
							},
						},
						"payload_type": {
							Description:  "Payload Type",
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: getNotificationPayloadTypes(),
							Default:      "exclude",
						},
						"paths": {
							Description: "Paths to Extract (If Payload Type is 'select_paths')",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							DefaultFunc: func() (interface{}, error) {
								return []string{}, nil
							},
						},
					},
				},
			},
		},
	}
}
