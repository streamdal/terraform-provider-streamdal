package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/streamdal/terraform-provider-streamdal/streamdal"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var types = []string{"replace", "delete", "truncate", "extract"}

func init() {
	schema.DescriptionKind = schema.StringMarkdown
	schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
		desc := s.Description
		if s.Default != nil {
			desc += fmt.Sprintf(" (Default: `%v`)", s.Default)
		}
		return strings.TrimSpace(desc)
	}
}

func New(version, apiToken string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"token": {
					Description: "Streamdal Server API token",
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("STREAMDAL_TOKEN", apiToken),
				},
				"address": {
					Description: "The address of the Streamdal server.",
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("STREAMDAL_ADDRESS", "localhost:8082"),
				},
				"connection_timeout": {
					Description: "The connection timeout for the Plumber server.",
					Type:        schema.TypeInt,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("STREAMDAL_CONNECTION_TIMEOUT", 10),
				},
			},
			ResourcesMap: map[string]*schema.Resource{
				"streamdal_pipeline":     resourcePipeline(),
				"streamdal_notification": resourceNotification(),
			},
			DataSourcesMap: map[string]*schema.Resource{
				"streamdal_pipeline":     dataSourcePipeline(),
				"streamdal_notification": dataSourceNotification(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(_ string, _ *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		cfg := &streamdal.Config{
			Address: d.Get("address").(string),
			Token:   d.Get("token").(string),
			Timeout: d.Get("connection_timeout").(int),
		}

		client, err := streamdal.New(cfg)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		return client, nil
	}
}

func dataSourceFiltersSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Description: "Field name to filter on",
					Type:        schema.TypeString,
					Required:    true,
				},

				"values": {
					Description: "Value(s) to filter by. Wildcards '*' are supported.",
					Type:        schema.TypeList,
					Required:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}

func buildFiltersDataSource(set *schema.Set) []*streamdal.Filter {
	var filters []*streamdal.Filter
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}
		filters = append(filters, &streamdal.Filter{
			Name:   m["name"].(string),
			Values: filterValues,
		})
	}
	return filters
}
