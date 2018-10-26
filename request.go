package postmanify

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Meetic/postmanify/postman2"
	"github.com/go-openapi/spec"
)

func (c *Converter) buildPostmanItem(url, method string, operation *spec.Operation) postman2.APIItem {

	//build request
	request := postman2.Request{
		Method: strings.ToUpper(method),
		URL:    c.buildPostmanURL(url, operation),
		Header: c.buildPostmanHeaders(operation),
	}

	if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
		request.Body = c.buildPostmanBody(operation)
	}

	//build item
	item := postman2.APIItem{
		Name:    url,
		Request: request,
	}

	if script := buildPostmanScript(operation.Extensions); len(script.Exec) > 0 {
		item.Event = []postman2.Event{
			postman2.Event{
				Listen: "test",
				Script: script,
			},
		}
	}

	return item

}

func (c *Converter) buildPostmanHeaders(operation *spec.Operation) []postman2.Header {
	headers := []postman2.Header{}
	if len(operation.Consumes) > 0 {
		if len(strings.TrimSpace(operation.Consumes[0])) > 0 {
			headers = append(headers, postman2.Header{
				Key:   "Content-Type",
				Value: strings.TrimSpace(operation.Consumes[0])})
		}
	}
	if len(operation.Produces) > 0 {
		if len(strings.TrimSpace(operation.Produces[0])) > 0 {
			headers = append(headers, postman2.Header{
				Key:   "Accept",
				Value: strings.TrimSpace(operation.Produces[0])})
		}
	}
	headers = append(headers, c.config.PostmanHeaders...)

	return headers
}

func (c *Converter) buildPostmanBody(operation *spec.Operation) postman2.RequestBody {

	requestBody := postman2.RequestBody{}

	var formData []postman2.FormData
	//formData
	for _, param := range operation.Parameters {

		if param.In == "formData" {
			var value string
			if param.Default != nil {
				value, _ = param.Default.(string)
			} else if param.Example != nil {
				value, _ = param.Example.(string)
			} else {
				value = "string"
			}

			formData = append(formData, postman2.FormData{
				Key:     param.Name,
				Value:   value,
				Enabled: param.Required,
				Type:    "text",
			})
		}

		//raw body
		if param.Required && param.In == "body" {
			if param.Schema.Type.Contains("object") {
				props := c.buildProperties(param.Schema.Properties)
				requestBody.Raw = props
			}

			if param.Schema.Type.Contains("array") {
				if param.Schema.Items.ContainsType("object") {
					array := []json.RawMessage{json.RawMessage(c.buildProperties(param.Schema.Items.Schema.Properties))}
					rawArray, _ := json.MarshalIndent(array, "", "\t")
					requestBody.Raw = string(rawArray)
					continue
				}

				var array []interface{}
				array = append(array, buildPropertyDefaultValue(param.Schema.Items.Schema.Type, param.Schema.Items.Schema.Format))
				rawArray, _ := json.MarshalIndent(array, "", "\t")
				requestBody.Raw = string(rawArray)

			}
		}
	}

	if len(formData) > 0 {
		requestBody.Mode = "formdata"
		requestBody.FormData = formData
		return requestBody
	}

	requestBody.Mode = "raw"
	return requestBody
}
