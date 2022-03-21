package util

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

func ParseKuskOpenAPI(yyml string) map[string]interface{} {
	var yml map[string]interface{}
	err := yaml.Unmarshal([]byte(yyml), &yml)
	if err != nil {
		fmt.Println(err)
	}
	fresh_yml := make(map[string]interface{})
	for k, v := range yml {
		switch v.(type) {
		case string:
			fresh_yml[k] = v
		case map[string]interface{}:
			if k != "x-kusk" {
				fresh_yml[k] = parseNode(k, v.(map[string]interface{}))
			}
		}

	}
	return fresh_yml
}

func parseNode(parent interface{}, v map[string]interface{}) map[string]interface{} {
	yml := make(map[string]interface{})

	for kk, vv := range v {
		switch vv.(type) {
		case string:
			if kk != "x-kusk" {
				yml[kk] = vv
			}
		case map[string]interface{}:
			if kk == "x-kusk" {
				ksk := getXKusk(vv)
				parentksk := getXKusk(parent)
				if parentksk != nil && ksk != nil {
					if parentksk.Disabled && !ksk.Disabled {
						yml[kk] = parseNode(vv, vv.(map[string]interface{}))
					}
				}
			}

			if kk != "x-kusk" {
				yml[kk] = parseNode(vv, vv.(map[string]interface{}))
			}
		case []interface{}:
			yml[kk] = vv
		default:
			break
		}
	}
	return yml
}

func getXKusk(value interface{}) *XKusk {
	jsn, _ := json.Marshal(value)
	xkusk := &XKusk{}

	if err := json.Unmarshal(jsn, xkusk); err != nil {
		fmt.Println(err)
	}
	return xkusk
}

type XKusk struct {
	Disabled bool `json:"disabled,omitempty"`
}
