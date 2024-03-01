package provider

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/streamdal/streamdal/libs/protos/build/go/protos"
	"github.com/streamdal/streamdal/libs/protos/build/go/protos/steps"
	"github.com/streamdal/terraform-provider-streamdal/streamdal"
)

func resourcePipeline() *schema.Resource {
	return &schema.Resource{
		Description: "Pipelines",

		CreateContext: resourcePipelineCreate,
		ReadContext:   resourcePipelineRead,
		UpdateContext: resourcePipelineUpdate,
		DeleteContext: resourcePipelineDelete,

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Pipeline ID",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"step": {
				Description: "Steps for this pipeline",
				Type:        schema.TypeList,
				Optional:    true,
				ConfigMode:  schema.SchemaConfigModeBlock,
				Elem:        getStepSchema(),
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}

}

// getConditionSchema returns the schema for a PipelineStepCondition message
// This is used to define the schema for the on_true, on_false, and on_error fields of a PipelineStep
func getConditionSchema() *schema.Resource {
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

// getStepSchema returns the schema for a PipelineStep message.
// This is in a separate method to try and keep the resourcePipeline method a bit cleaner
func getStepSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Step Name",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"on_true": {
				Description: "Determines the next action if the result of the step is true",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        getConditionSchema(),
			},
			"on_false": {
				Description: "Determines the next action if the result of the step is false",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        getConditionSchema(),
			},
			"on_error": {
				Description: "Determines the next action if the result of the step is an error",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        getConditionSchema(),
			},
			"dynamic": {
				Description: "Should this step use the result from the previous step",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"detective": {
				Description: "Detective Step",
				Type:        schema.TypeList,
				Optional:    true,
				Required:    false,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Description: "Path",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"type": {
							Description:  "Detective Type",
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: getDetectiveTypes(),
						},
						"args": {
							Description: "Arguments",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							DefaultFunc: func() (interface{}, error) {
								return []string{}, nil
							},
						},
						"negate": {
							Description: "Negate",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
					},
				},
			},
			"transform": {
				Description: "Transform Step",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description:  "Transform Type",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: getTransformTypes(),
						},
						"replace_value": {
							Description: "Replace value of a field",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path": {
										Description: "Path",
										Type:        schema.TypeString,
										Optional:    true, // TODO: make optional or required based on value of "dynamic"
									},
									"value": {
										Description: "Value",
										Type:        schema.TypeString,
										Optional:    true, // TODO: make optional or required based on value of "dynamic"
									},
								},
							},
						},
						"delete_field": {
							Description: "Delete field",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"paths": {
										Description: "Paths",
										Type:        schema.TypeList,
										Optional:    true,
										DefaultFunc: func() (interface{}, error) {
											return []string{}, nil
										},
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"obfuscate": {
							Description: "Obfuscate value",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path": {
										Description: "Path",
										Type:        schema.TypeString,
										Optional:    true,
									},
								},
							},
						},
						"mask": {
							Description: "Mask value",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path": {
										Description: "Path",
										Type:        schema.TypeString,
										Optional:    true,
									},
									"mask": {
										Description: "Mask",
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "*",
									},
								},
							},
						},

						"truncate": {
							Description: "Truncate value",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Description:  "Truncate Type",
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: getTransformTruncateTypes(),
									},
									"path": {
										Description: "Path",
										Type:        schema.TypeString,
										Optional:    true,
									},
								},
							},
						},
						"extract": {
							Description: "Extract value",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"paths": {
										Description: "Paths",
										Type:        schema.TypeList,
										Optional:    true,
										DefaultFunc: func() (interface{}, error) {
											return []string{}, nil
										},
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
			"http_request": {
				Description: "HTTP Request Step",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"method": {
							Description:  "HTTP Method",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: getHttpMethods(),
						},
						"url": {
							Description: "URL",
							Type:        schema.TypeString,
							Required:    true,
						},
						"headers": {
							Description: "Headers",
							Type:        schema.TypeMap,
							Optional:    true,
							DefaultFunc: func() (interface{}, error) {
								return map[string]interface{}{}, nil
							},
						},
						"body": {
							Description: "Body",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
			"valid_json": {
				Description: "Valid JSON Step",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        &schema.Resource{},
			},
			"schema_validation": {
				Description: "Schema Validation Step",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description:  "Schema Validation Type",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: getSchemaValidationTypes(),
						},
						"condition": {
							Description:  "Schema Validation Condition",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: getSchemaValidationConditions(),
						},
						"json_schema": {
							Description: "JSON Schema",
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"draft": {
										Description:  "JSON Schema Draft",
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: getSchemaValidationJSONSchemaDrafts(),
									},
									"json_schema": {
										Description: "Schema Definition",
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

func resourcePipelineRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	s := m.(*streamdal.Streamdal)

	resp, err := s.GetPipeline(ctx, &protos.GetPipelineRequest{
		PipelineId: d.Id(),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	opts := resp.GetPipeline()

	d.SetId(opts.GetId())
	_ = d.Set("name", opts.GetName())
	_ = d.Set("steps", opts.GetSteps())

	return diags
}

func resourcePipelineCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	pc := m.(*streamdal.Streamdal)

	pipeline, moreDiags := buildPipeline(d)
	if moreDiags.HasError() {
		return append(diags, moreDiags...)
	}

	resp, err := pc.CreatePipeline(ctx, &protos.CreatePipelineRequest{
		Pipeline: pipeline,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.PipelineId)

	return diags
}

func resourcePipelineUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var p *protos.Pipeline

	p, moreDiags := buildPipeline(d)
	if moreDiags.HasError() {
		return append(diags, moreDiags...)
	}

	client := m.(*streamdal.Streamdal)
	_, err := client.UpdatePipeline(ctx, &protos.UpdatePipelineRequest{
		Pipeline: p,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourcePipelineDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	s := m.(*streamdal.Streamdal)
	_, err := s.DeletePipeline(ctx, &protos.DeletePipelineRequest{
		PipelineId: d.Id(),
	})

	if err != nil {
		return diag.Errorf("Error deleting pipeline: %s", s.Token)
		//return diag.FromErr(err)
	}

	return diags
}

func buildPipeline(d *schema.ResourceData) (*protos.Pipeline, diag.Diagnostics) {
	var diags diag.Diagnostics
	p := &protos.Pipeline{
		Name:    d.Get("name").(string),
		Steps:   []*protos.PipelineStep{},
		XPaused: proto.Bool(false),
	}

	pipelineSteps := d.Get("step").([]interface{})

	for _, step := range pipelineSteps {
		stepMap := step.(map[string]interface{})

		//return nil, diag.Errorf("%#v", stepMap)

		onTrue, diags := generateCondition(stepMap, "on_true")
		if diags.HasError() {
			return nil, diags
		}

		onFalse, diags := generateCondition(stepMap, "on_false")
		if diags.HasError() {
			return nil, diags
		}

		onError, diags := generateCondition(stepMap, "on_error")
		if diags.HasError() {
			return nil, diags
		}

		s := &protos.PipelineStep{
			Name:    stepMap["name"].(string),
			OnTrue:  onTrue,
			OnFalse: onFalse,
			OnError: onError,
			Dynamic: stepMap["dynamic"].(bool),
		}

		t := getStepType(stepMap)

		generateStep(s, stepMap, t)

		p.Steps = append(p.Steps, s)

	}

	return p, diags
}

func generateCondition(stepMap map[string]interface{}, conditionType string) (*protos.PipelineStepConditions, diag.Diagnostics) {
	var diags diag.Diagnostics

	onCondition, ok := stepMap[conditionType].([]interface{})
	if !ok || len(onCondition) != 1 {
		// Condition not specified, value will be nil in PipelineStep
		return nil, diags
	}

	conditionCfg, ok := onCondition[0].(map[string]interface{})
	if !ok {
		return nil, diag.Errorf("Error generating condition: %s", "Invalid condition config")
	}

	acStr := conditionCfg["abort"].(string)
	acType, err := abortConditionFromString(acStr)
	if err != nil {
		return nil, diag.Errorf("Error generating abort condition: %s", err)
	}

	cond := &protos.PipelineStepConditions{
		Abort:    acType,
		Metadata: interfaceMapToStringMap(conditionCfg["metadata"]),
	}

	if notCfg, ok := conditionCfg["notification"].([]interface{}); ok && len(notCfg) > 0 {
		payloadTypeStr := notCfg[0].(map[string]interface{})["payload_type"].(string)

		payloadType, err := notificationPayloadTypeFromString(payloadTypeStr)
		if err != nil {
			return nil, diag.Errorf("Error generating notification payload type: %s", err)
		}

		cond.Notification = &protos.PipelineStepNotification{
			NotificationConfigIds: interfaceToStrings(conditionCfg["notification_config_ids"]),
			PayloadType:           payloadType,
			Paths:                 interfaceToStrings(conditionCfg["paths"]),
		}
	}

	return cond, diags
}

func generateStep(s *protos.PipelineStep, stepMap map[string]interface{}, t string) diag.Diagnostics {
	var diags diag.Diagnostics

	switch t {
	case "detective":
		diags = generateStepDetective(s, stepMap)
		if diags.HasError() {
			diags = append(diags, diags...)
		}
	case "transform":
		diags = generateStepTransform(s, stepMap)
		if diags.HasError() {
			diags = append(diags, diags...)
		}
	case "http_request":
		diags = generateStepHttpRequest(s, stepMap)
		if diags.HasError() {
			diags = append(diags, diags...)
		}
	case "valid_json":
		diags = generateValidJsonStep(s, stepMap)
		if diags.HasError() {
			diags = append(diags, diags...)
		}
	case "schema_validation":
		diags = generateSchemaValidationStep(s, stepMap)
		if diags.HasError() {
			diags = append(diags, diags...)
		}
	case "kv":
		diags = generateKVStep(s, stepMap)
		if diags.HasError() {
			diags = append(diags, diags...)

		}
	default:
		return diag.Errorf("Unknown step type: %s", t)
	}

	return diags
}

func generateKVStep(s *protos.PipelineStep, stepMap map[string]interface{}) diag.Diagnostics {
	stepData := stepMap["kv"].([]interface{})
	config := stepData[0].(map[string]interface{})

	mode, err := kvModeFromString(config["type"].(string))
	if err != nil {
		return diag.Errorf("Error generating kv step: %s", err)
	}

	action, err := kvActionFromString(config["action"].(string))
	if err != nil {
		return diag.Errorf("Error generating kv step: %s", err)
	}

	s.Step = &protos.PipelineStep_Kv{
		Kv: &steps.KVStep{
			Mode:   mode,
			Action: action,
			Key:    config["key"].(string),
			Value:  []byte(config["value"].(string)),
		},
	}

	return diag.Diagnostics{}

}

func generateSchemaValidationStep(s *protos.PipelineStep, stepMap map[string]interface{}) diag.Diagnostics {
	stepData := stepMap["schema_validation"].([]interface{})
	config := stepData[0].(map[string]interface{})

	t, err := schemaValidationTypeFromString(config["type"].(string))
	if err != nil {
		return diag.Errorf("Error generating schema validation step: %s", err)

	}

	cond, err := schemaValidationConditionFromString(config["condition"].(string))
	if err != nil {
		return diag.Errorf("Error generating schema validation step: %s", err)
	}

	step := &protos.PipelineStep_SchemaValidation{
		SchemaValidation: &steps.SchemaValidationStep{
			Type:      t,
			Condition: cond,
			Options:   nil, // Filled out below
		},
	}

	switch t {
	case steps.SchemaValidationType_SCHEMA_VALIDATION_TYPE_JSONSCHEMA:
		draft, err := schemaValidationJSONSchemaDraftFromString(config["draft"].(string))
		if err != nil {
			return diag.Errorf("Error generating schema validation step: %s", err)
		}
		step.SchemaValidation.Options = &steps.SchemaValidationStep_JsonSchema{
			JsonSchema: &steps.SchemaValidationJSONSchema{
				JsonSchema: []byte(config["schema"].(string)),
				Draft:      draft,
			},
		}
	default:
		return diag.Errorf("Error generating schema validation step: unknown schema validation type: %s", t)
	}

	s.Step = step

	return diag.Diagnostics{}

}

func generateValidJsonStep(s *protos.PipelineStep, stepMap map[string]interface{}) diag.Diagnostics {
	s.Step = &protos.PipelineStep_ValidJson{
		ValidJson: &steps.ValidJSONStep{},
	}

	return diag.Diagnostics{}
}

func generateStepHttpRequest(s *protos.PipelineStep, stepMap map[string]interface{}) diag.Diagnostics {
	stepData := stepMap["http_request"].([]interface{})
	config := stepData[0].(map[string]interface{})

	t, err := httpMethodFromString(config["method"].(string))
	if err != nil {
		return diag.Errorf("Error generating http request step: %s", err)
	}

	s.Step = &protos.PipelineStep_HttpRequest{
		HttpRequest: &steps.HttpRequestStep{
			Request: &steps.HttpRequest{
				Method:  t,
				Url:     config["url"].(string),
				Headers: interfaceMapToStringMap(config["headers"]),
				Body:    []byte(config["body"].(string)),
			},
		},
	}

	return diag.Diagnostics{}

}

func generateStepDetective(s *protos.PipelineStep, stepMap map[string]interface{}) diag.Diagnostics {
	stepData := stepMap["detective"].([]interface{})
	config := stepData[0].(map[string]interface{})

	t, err := detectiveTypeFromString(config["type"].(string))
	if err != nil {
		return diag.Errorf("Error generating detective step: %s", err)
	}

	s.Step = &protos.PipelineStep_Detective{
		Detective: &steps.DetectiveStep{
			Path:   proto.String(config["path"].(string)),
			Args:   interfaceToStrings(config["args"]),
			Negate: proto.Bool(config["negate"].(bool)),
			Type:   t,
		},
	}

	return diag.Diagnostics{}
}

func generateStepTransform(s *protos.PipelineStep, stepMap map[string]interface{}) diag.Diagnostics {
	stepData := stepMap["transform"].([]interface{})
	config := stepData[0].(map[string]interface{})

	t, err := transformTypeFromString(config["type"].(string))
	if err != nil {
		return diag.Errorf("Error generating transform step: %s", err)
	}

	s.Step = &protos.PipelineStep_Transform{
		Transform: &steps.TransformStep{
			Type: t,
		},
	}

	// Populate transform oneof
	switch t {
	case steps.TransformType_TRANSFORM_TYPE_REPLACE_VALUE:
		replaceData, ok := config["replace_value"].([]interface{})
		if !ok || len(replaceData) == 0 {
			return diag.Errorf("Error generating transform step: replace value config not found")
		}

		replaceCfg := replaceData[0].(map[string]interface{})

		s.GetTransform().Options = &steps.TransformStep_ReplaceValueOptions{
			ReplaceValueOptions: &steps.TransformReplaceValueOptions{
				Path:  replaceCfg["path"].(string),
				Value: replaceCfg["value"].(string),
			},
		}
	case steps.TransformType_TRANSFORM_TYPE_DELETE_FIELD:
		deleteData, ok := config["delete_field"].([]interface{})
		if !ok || len(deleteData) == 0 {
			return diag.Errorf("Error generating transform step: delete field config not found")
		}

		deleteCfg := deleteData[0].(map[string]interface{})

		s.GetTransform().Options = &steps.TransformStep_DeleteFieldOptions{
			DeleteFieldOptions: &steps.TransformDeleteFieldOptions{
				Paths: interfaceToStrings(deleteCfg["path"]),
			},
		}
	case steps.TransformType_TRANSFORM_TYPE_OBFUSCATE_VALUE:
		obfuscateData, ok := config["obfuscate"].([]interface{})
		if !ok || len(obfuscateData) == 0 {
			return diag.Errorf("Error generating transform step: obfuscate value config not found")
		}

		obfuscateCfg := obfuscateData[0].(map[string]interface{})

		s.GetTransform().Options = &steps.TransformStep_ObfuscateOptions{
			ObfuscateOptions: &steps.TransformObfuscateOptions{
				Path: obfuscateCfg["path"].(string),
			},
		}
	case steps.TransformType_TRANSFORM_TYPE_MASK_VALUE:
		maskData, ok := config["mask"].([]interface{})
		if !ok || len(maskData) == 0 {
			return diag.Errorf("Error generating transform mask step: mask value config not found")
		}

		maskCfg := maskData[0].(map[string]interface{})

		s.GetTransform().Options = &steps.TransformStep_MaskOptions{
			MaskOptions: &steps.TransformMaskOptions{
				Path: maskCfg["path"].(string),
				Mask: maskCfg["mask"].(string),
			},
		}
	case steps.TransformType_TRANSFORM_TYPE_TRUNCATE_VALUE:
		truncateData, ok := config["truncate"].([]interface{})
		if !ok || len(truncateData) == 0 {
			return diag.Errorf("Error generating transform step: truncate value config not found")
		}

		truncateCfg := truncateData[0].(map[string]interface{})

		tt, err := transformTruncateTypeFromString(truncateCfg["type"].(string))
		if err != nil {
			return diag.Errorf("Error generating transform truncate step: %s", err)
		}

		s.GetTransform().Options = &steps.TransformStep_TruncateOptions{
			TruncateOptions: &steps.TransformTruncateOptions{
				Type:  tt,
				Path:  truncateCfg["path"].(string),
				Value: truncateCfg["type"].(int32),
			},
		}
	case steps.TransformType_TRANSFORM_TYPE_EXTRACT:
		extractData, ok := config["extract"].([]interface{})
		if !ok || len(extractData) == 0 {
			return diag.Errorf("Error generating transform step: extract value config not found")
		}

		extractCfg := extractData[0].(map[string]interface{})

		s.GetTransform().Options = &steps.TransformStep_ExtractOptions{
			ExtractOptions: &steps.TransformExtractOptions{
				Paths:   interfaceToStrings(extractCfg["paths"]),
				Flatten: extractCfg["flatten"].(bool),
			},
		}
	default:
		return diag.Errorf("Error generating transform step: unknown transform type: %s", t)
	}

	return diag.Diagnostics{}
}
