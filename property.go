package postmanify

import (
	"encoding/json"
	"sort"
	"time"

	"github.com/go-openapi/spec"
)

func (c *Converter) buildProperties(properties map[string]spec.Schema) string {

	body := make(map[string]interface{})

	keys := []string{}
	for key := range properties {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		prop := properties[key]

		//Property as an example value : we take it as value
		if prop.Example != nil {
			body[key] = prop.Example
			continue
		}

		//Property as a Enum : we take the first possible value
		//Note: we only support string enum for now;
		//TODO : add support for other type enum.
		if prop.Type.Contains("string") && len(prop.Enum) > 0 {
			body[key] = prop.Enum[0]
			continue
		}

		if prop.Type.Contains("object") {
			body[key] = json.RawMessage(c.buildProperties(prop.Properties))
			continue
		}

		if prop.Type.Contains("array") {
			if prop.Items.ContainsType("object") {
				array := []json.RawMessage{json.RawMessage(c.buildProperties(prop.Items.Schema.Properties))}
				rawArray, _ := json.MarshalIndent(array, "", "\t")
				body[key] = json.RawMessage(rawArray)
				continue
			}
			var array []interface{}
			array = append(array, buildPropertyDefaultValue(prop.Items.Schema.Type, prop.Items.Schema.Format))
			rawArray, _ := json.MarshalIndent(array, "", "\t")
			body[key] = json.RawMessage(rawArray)
			continue
		}

		body[key] = buildPropertyDefaultValue(prop.Type, prop.Format)

	}

	b, err := json.MarshalIndent(body, "", "\t")
	if err != nil {
		panic(err)
	}

	return string(b)
}

func buildPropertyDefaultValue(propType spec.StringOrArray, propFormat string) interface{} {

	if propType == nil {
		return ""
	}

	//Property has no example value : we set one by default
	if propType.Contains("integer") {
		return 0
	}

	if propType.Contains("string") {
		switch propFormat {
		case "date-time":
			return time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC).Format(time.RFC3339)
		default:
			return "string"
		}
	}

	return ""
}
