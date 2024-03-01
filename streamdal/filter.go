package streamdal

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/minio/pkg/wildcard"
)

type Filter struct {
	Name   string
	Values []string
}

// filterJSON takes a slice of map[string]interface{}, which represents the output of an gRPC Call as JSON.
// It then applies the given filters to the JSON array and returns only entries in the data param which
// have a field that matches a filter.
func filterJSON(data map[string]interface{}, filters []*Filter) ([]map[string]interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	found := make([]map[string]interface{}, 0)

	for _, filter := range filters {
		for _, item := range data {
			row := item.(map[string]interface{})
			// Can't match on non-existent keys
			if _, ok := row[filter.Name]; !ok {
				return nil, diag.FromErr(fmt.Errorf("%s is not a valid key", filter.Name))
			}

			checkVal := fmt.Sprintf("%s", row[filter.Name])
			for _, val := range filter.Values {
				// Wildcard match
				if matches(val, checkVal) {
					found = append(found, row)
				}
			}
		}
	}

	return found, diags
}

func matches(val, checkVal string) bool {
	if strings.Contains(val, "*") && wildcard.MatchSimple(val, checkVal) {
		return true
	}

	if val == checkVal {
		return true
	}

	return false
}

func (f *Filter) String() string {
	return fmt.Sprintf("%s: %s", f.Name, strings.Join(f.Values, ","))
}

func filterString(filters []*Filter) string {
	out := ""
	for _, filter := range filters {
		out += filter.String() + "\n"
	}
	return out
}
