package postmanify

import (
	"encoding/json"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/Meetic/postmanify/postman2"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
)

type Config struct {
	Hostname       string
	HostnamePrefix string
	HostnameSuffix string
	Schema         string
	BasePath       string
	PostmanHeaders []postman2.Header
}

type Converter struct {
	config Config
}

func NewConverter(cfg Config) *Converter {
	return &Converter{
		config: cfg,
	}
}

func (c *Converter) Convert(swaggerSpec []byte) ([]byte, error) {

	specDoc, err := loads.Analyzed(swaggerSpec, "2.0")
	if err != nil {
		return nil, err
	}

	specDocExpand, err := specDoc.Expanded(&spec.ExpandOptions{
		SkipSchemas:         false,
		ContinueOnError:     true,
		AbsoluteCircularRef: true,
	})
	if err != nil {
		return nil, err
	}

	swag := specDocExpand.Spec()

	if c.config.Hostname == "" {
		c.config.Hostname = strings.TrimSpace(swag.Host)
	}

	if c.config.BasePath == "" {
		c.config.BasePath = strings.TrimSpace(swag.BasePath)
	}

	//if schema is not defined in config, we take the first one declared on the swagger specs.
	if c.config.Schema == "" && len(swag.Schemes) >= 1 {
		c.config.Schema = strings.TrimSpace(swag.Schemes[0])
	} else {
		c.config.Schema = "http"
	}

	pman := postman2.NewCollection(strings.TrimSpace(swag.Info.Title), strings.TrimSpace(swag.Info.Description))

	if err := c.addUrls(swag.Paths.Paths, &pman); err != nil {
		return nil, err
	}

	return json.MarshalIndent(pman, "", "  ")

}

func (c *Converter) addUrls(paths map[string]spec.PathItem, pman *postman2.Collection) error {
	urls := []string{}
	for url := range paths {
		urls = append(urls, url)
	}
	sort.Strings(urls)

	for _, url := range urls {
		path := paths[url]

		if pathHasMethodWithTag(path, http.MethodGet) {
			pman.AddItem(c.buildPostmanItem(url, http.MethodGet, path.Get), strings.TrimSpace(path.Get.Tags[0]))
		}
		if pathHasMethodWithTag(path, http.MethodPatch) {
			pman.AddItem(c.buildPostmanItem(url, http.MethodPatch, path.Patch), strings.TrimSpace(path.Patch.Tags[0]))
		}
		if pathHasMethodWithTag(path, http.MethodPost) {
			pman.AddItem(c.buildPostmanItem(url, http.MethodPost, path.Post), strings.TrimSpace(path.Post.Tags[0]))
		}
		if pathHasMethodWithTag(path, http.MethodPut) {
			pman.AddItem(c.buildPostmanItem(url, http.MethodPut, path.Put), strings.TrimSpace(path.Put.Tags[0]))
		}
		if pathHasMethodWithTag(path, http.MethodDelete) {
			pman.AddItem(c.buildPostmanItem(url, http.MethodDelete, path.Delete), strings.TrimSpace(path.Delete.Tags[0]))
		}
	}

	return nil
}

func pathHasMethodWithTag(path spec.PathItem, method string) bool {
	method = strings.TrimSpace(strings.ToLower(method))
	switch method {
	case "get":
		if path.Get != nil && len(path.Get.Tags) > 0 && len(strings.TrimSpace(path.Get.Tags[0])) > 0 {
			return true
		}
	case "patch":
		if path.Patch != nil && len(path.Patch.Tags) > 0 && len(strings.TrimSpace(path.Patch.Tags[0])) > 0 {
			return true
		}
	case "post":
		if path.Post != nil && len(path.Post.Tags) > 0 && len(strings.TrimSpace(path.Post.Tags[0])) > 0 {
			return true
		}
	case "put":
		if path.Put != nil && len(path.Put.Tags) > 0 && len(strings.TrimSpace(path.Put.Tags[0])) > 0 {
			return true
		}
	case "delete":
		if path.Delete != nil && len(path.Delete.Tags) > 0 && len(strings.TrimSpace(path.Delete.Tags[0])) > 0 {
			return true
		}
	}
	return false
}

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
	return postman2.APIItem{
		Name:    url,
		Request: request,
	}

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

func (c *Converter) buildPostmanURL(url string, operation *spec.Operation) postman2.URL {

	//create hostname
	host := strings.TrimSpace(strings.Join([]string{
		c.config.HostnamePrefix,
		c.config.Hostname,
		c.config.HostnameSuffix,
	}, ""))

	//create URI
	rawPostmanURL := strings.TrimSpace(strings.Join([]string{
		host,
		c.config.BasePath,
		strings.Replace(url, " ", "", -1),
	}, "/"))

	rx1 := regexp.MustCompile(`/+`)
	rawPostmanURL = rx1.ReplaceAllString(rawPostmanURL, "/")
	rx2 := regexp.MustCompile(`^/+`)
	rawPostmanURL = rx2.ReplaceAllString(rawPostmanURL, "")

	//Add schema
	rawPostmanURL = strings.Join([]string{c.config.Schema, rawPostmanURL}, "://")

	//Replace URI parameters
	rx3 := regexp.MustCompile(`(^|[^\{])\{([^\/\{\}]+)\}([^\}]|$)`)
	rawPostmanURL = rx3.ReplaceAllString(rawPostmanURL, "$1{{$2}}$3")

	postmanURL := postman2.NewURL(rawPostmanURL)

	// Set Default URL Path Parameters
	rx4 := regexp.MustCompile(`^\s*({{(.+)}})\s*$`)

	for _, part := range postmanURL.Path {
		rs4 := rx4.FindAllStringSubmatch(part, -1)
		if len(rs4) > 0 {
			baseVariable := rs4[0][2]
			var defaultValue interface{}
			for _, parameter := range operation.Parameters {
				if parameter.Name == baseVariable {
					defaultValue = parameter.Default
					break
				}
			}
			postmanURL.AddVariable(baseVariable, defaultValue)
		}
	}

	return postmanURL
}

func (c *Converter) buildPostmanBody(operation *spec.Operation) postman2.RequestBody {

	requestBody := postman2.RequestBody{
		Mode: "raw",
	}

	for _, param := range operation.Parameters {
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

	return requestBody
}

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
