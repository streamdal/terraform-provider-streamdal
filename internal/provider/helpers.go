package provider

import (
	"errors"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/streamdal/streamdal/libs/protos/build/go/protos"
	"github.com/streamdal/streamdal/libs/protos/build/go/protos/shared"
	"github.com/streamdal/streamdal/libs/protos/build/go/protos/steps"
)

// interfaceToStrings converts an interface{} value to []string
// This is needed when a nested resource, say "kafka" has a value that is a schema.TypeList of schema.TypeString
func interfaceToStrings(value interface{}) []string {
	strs := make([]string, 0)

	if value == nil {
		return strs
	}

	for _, v := range value.([]interface{}) {
		strs = append(strs, v.(string))
	}

	return strs
}

func interfaceMapToStringMap(value interface{}) map[string]string {
	m := make(map[string]string)

	if value == nil {
		return m
	}

	for k, v := range value.(map[string]interface{}) {
		m[k] = v.(string)
	}

	return m
}

var stepTypes = []string{"detective", "transform", "http_request", "valid_json", "schema_validation", "kv"}

func getStepType(d map[string]interface{}) string {
	if d == nil {
		return ""
	}

	for _, st := range stepTypes {
		if opts, ok := d[st].([]interface{}); ok && len(opts) > 0 {
			return st
		}
	}

	return ""
}

// detectiveTypeFromString converts a string to a detective type enum
func detectiveTypeFromString(s string) (steps.DetectiveType, error) {
	for id, v := range steps.DetectiveType_name {
		v = strings.Replace(v, "DETECTIVE_TYPE_", "", -1)
		v = strings.ToLower(v)
		if s == v {
			return steps.DetectiveType(id), nil
		}
	}

	return 0, errors.New("invalid detective type")
}

func transformTypeFromString(s string) (steps.TransformType, error) {
	for id, v := range steps.TransformType_name {
		v = strings.Replace(v, "TRANSFORM_TYPE_", "", -1)
		v = strings.ToLower(v)
		if s == v {
			return steps.TransformType(id), nil
		}
	}

	return 0, errors.New("invalid transform type")

}

// getDetectiveTypes returns all detective type enums as a slice of strings
func getDetectiveTypes() schema.SchemaValidateFunc {
	t := make([]string, 0)

	for _, v := range steps.DetectiveType_name {
		v = strings.Replace(v, "DETECTIVE_TYPE_", "", -1)
		v = strings.ToLower(v)
		t = append(t, v)
	}

	return validation.StringInSlice(t, true)
}

// getDetectiveTypes returns all transform type enums as a slice of strings
func getTransformTypes() schema.SchemaValidateFunc {
	t := make([]string, 0)

	for _, v := range steps.TransformType_name {
		v = strings.Replace(v, "TRANSFORM_TYPE_", "", -1)
		v = strings.ToLower(v)
		t = append(t, v)
	}

	return validation.StringInSlice(t, true)
}

func getTransformTruncateTypes() schema.SchemaValidateFunc {
	t := make([]string, 0)

	for _, v := range steps.TransformTruncateType_name {
		v = strings.Replace(v, "TRANSFORM_TRUNCATE_TYPE_", "", -1)
		v = strings.ToLower(v)
		t = append(t, v)
	}

	return validation.StringInSlice(t, true)
}

func transformTruncateTypeFromString(s string) (steps.TransformTruncateType, error) {
	for id, v := range steps.TransformTruncateType_name {
		v = strings.Replace(v, "TRANSFORM_TRUNCATE_TYPE_", "", -1)
		v = strings.ToLower(v)
		if s == v {
			return steps.TransformTruncateType(id), nil
		}
	}

	return 0, errors.New("invalid transform truncate type")
}

func getAbortConditions() schema.SchemaValidateFunc {
	t := make([]string, 0)

	for _, v := range protos.AbortCondition_name {
		v = strings.Replace(v, "ABORT_CONDITION_", "", -1)
		v = strings.ToLower(v)
		t = append(t, v)
	}

	return validation.StringInSlice(t, true)
}

func abortConditionFromString(s string) (protos.AbortCondition, error) {
	for id, v := range protos.AbortCondition_name {
		v = strings.Replace(v, "ABORT_CONDITION_", "", -1)
		v = strings.ToLower(v)
		if s == v {
			return protos.AbortCondition(id), nil
		}
	}

	return 0, errors.New("invalid abort condition")

}

func getNotificationPayloadTypes() schema.SchemaValidateFunc {
	t := make([]string, 0)

	for _, v := range protos.PipelineStepNotification_PayloadType_name {
		v = strings.Replace(v, "PAYLOAD_TYPE_", "", -1)
		v = strings.ToLower(v)
		t = append(t, v)
	}

	return validation.StringInSlice(t, true)
}

func notificationPayloadTypeFromString(s string) (protos.PipelineStepNotification_PayloadType, error) {
	for id, v := range protos.PipelineStepNotification_PayloadType_name {
		v = strings.Replace(v, "PAYLOAD_TYPE_", "", -1)
		v = strings.ToLower(v)
		if s == v {
			return protos.PipelineStepNotification_PayloadType(id), nil
		}
	}

	return 0, errors.New("invalid notification payload type")
}

func getHttpMethods() schema.SchemaValidateFunc {
	t := make([]string, 0)

	for _, v := range steps.HttpRequestMethod_name {
		v = strings.Replace(v, "HTTP_REQUEST_METHOD_", "", -1)
		v = strings.ToLower(v)
		t = append(t, v)
	}

	return validation.StringInSlice(t, true)
}

func httpMethodFromString(s string) (steps.HttpRequestMethod, error) {
	for id, v := range steps.HttpRequestMethod_name {
		v = strings.Replace(v, "HTTP_REQUEST_METHOD_", "", -1)
		v = strings.ToLower(v)
		if s == v {
			return steps.HttpRequestMethod(id), nil
		}
	}

	return 0, errors.New("invalid http method")
}

func getSchemaValidationTypes() schema.SchemaValidateFunc {
	t := make([]string, 0)

	for _, v := range steps.SchemaValidationType_name {
		v = strings.Replace(v, "SCHEMA_VALIDATION_TYPE_", "", -1)
		v = strings.ToLower(v)
		t = append(t, v)
	}

	return validation.StringInSlice(t, true)
}

func schemaValidationTypeFromString(s string) (steps.SchemaValidationType, error) {
	for id, v := range steps.SchemaValidationType_name {
		v = strings.Replace(v, "SCHEMA_VALIDATION_TYPE_", "", -1)
		v = strings.ToLower(v)
		if s == v {
			return steps.SchemaValidationType(id), nil
		}
	}

	return 0, errors.New("invalid schema validation type")
}

func getSchemaValidationConditions() schema.SchemaValidateFunc {
	t := make([]string, 0)

	for _, v := range steps.SchemaValidationCondition_name {
		v = strings.Replace(v, "SCHEMA_VALIDATION_CONDITION_", "", -1)
		v = strings.ToLower(v)
		t = append(t, v)
	}

	return validation.StringInSlice(t, true)
}

func schemaValidationConditionFromString(s string) (steps.SchemaValidationCondition, error) {
	for id, v := range steps.SchemaValidationCondition_name {
		v = strings.Replace(v, "SCHEMA_VALIDATION_CONDITION_", "", -1)
		v = strings.ToLower(v)
		if s == v {
			return steps.SchemaValidationCondition(id), nil
		}
	}

	return 0, errors.New("invalid schema validation condition")

}

func getSchemaValidationJSONSchemaDrafts() schema.SchemaValidateFunc {
	t := make([]string, 0)

	for _, v := range steps.JSONSchemaDraft_name {
		v = strings.Replace(v, "JSON_SCHEMA_", "", -1)
		v = strings.ToLower(v)
		t = append(t, v)
	}

	return validation.StringInSlice(t, true)
}

func schemaValidationJSONSchemaDraftFromString(s string) (steps.JSONSchemaDraft, error) {
	for id, v := range steps.JSONSchemaDraft_name {
		v = strings.Replace(v, "JSON_SCHEMA_", "", -1)
		v = strings.ToLower(v)
		if s == v {
			return steps.JSONSchemaDraft(id), nil
		}
	}

	return 0, errors.New("invalid schema validation JSON schema draft")
}

func getKvTypes() schema.SchemaValidateFunc {
	t := make([]string, 0)

	for _, v := range steps.KVMode_name {
		v = strings.Replace(v, "KV_MODE_", "", -1)
		v = strings.ToLower(v)
		t = append(t, v)
	}

	return validation.StringInSlice(t, true)
}

func kvModeFromString(s string) (steps.KVMode, error) {
	for id, v := range steps.KVMode_name {
		v = strings.Replace(v, "KV_MODE_", "", -1)
		v = strings.ToLower(v)
		if s == v {
			return steps.KVMode(id), nil
		}
	}

	return 0, errors.New("invalid kv mode")
}

func getKvActions() schema.SchemaValidateFunc {
	t := make([]string, 0)

	for _, v := range shared.KVAction_name {
		v = strings.Replace(v, "KV_ACTION_", "", -1)
		v = strings.ToLower(v)
		t = append(t, v)
	}

	return validation.StringInSlice(t, true)
}

func kvActionFromString(s string) (shared.KVAction, error) {
	for id, v := range shared.KVAction_name {
		v = strings.Replace(v, "KV_ACTION_", "", -1)
		v = strings.ToLower(v)
		if s == v {
			return shared.KVAction(id), nil
		}
	}

	return 0, errors.New("invalid kv action")
}
