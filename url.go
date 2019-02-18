package postmanify

import (
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/seblegall/postmanify/postman2"
	"github.com/go-openapi/spec"
)

//addUrls add a postman items for each swagger path in the spec
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

//buildPostmanURL build a postman url, part of a postman item, from a swagger operation
func (c *Converter) buildPostmanURL(url string, operation *spec.Operation) postman2.URL {

	//create hostname
	host := strings.TrimSpace(strings.Join([]string{
		c.config.Hostname,
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

	queryParams := buildQueryParams(operation)

	for _, queryParam := range queryParams {
		postmanURL.AddQueryParam(queryParam)
	}

	return postmanURL
}

//buildQueryParams build postman query param from a swagger operation spec
func buildQueryParams(operation *spec.Operation) []postman2.URLQueryParam {

	var queryParam []postman2.URLQueryParam

	for _, param := range operation.Parameters {
		if param.In == "query" {

			if param.Example != nil {
				queryParam = append(queryParam, postman2.URLQueryParam{Key: param.Name, Value: param.Example})
				continue
			}

			if param.Default != nil {
				queryParam = append(queryParam, postman2.URLQueryParam{Key: param.Name, Value: param.Default})
				continue
			}

			if len(param.Enum) > 0 {
				queryParam = append(queryParam, postman2.URLQueryParam{Key: param.Name, Value: param.Enum[0]})
				continue
			}

			if param.Type == "array" {
				if param.Items.Example != nil {
					queryParam = append(queryParam, postman2.URLQueryParam{Key: param.Name, Value: param.Items.Example})
					continue
				}
				if param.Items.Default != nil {
					queryParam = append(queryParam, postman2.URLQueryParam{Key: param.Name, Value: param.Items.Default})
					continue
				}

				if len(param.Items.Enum) > 0 {
					queryParam = append(queryParam, postman2.URLQueryParam{Key: param.Name, Value: param.Items.Enum[0]})
					continue
				}

				queryParam = append(queryParam, postman2.URLQueryParam{Key: param.Name, Value: buildQueryParamDefaultValue(param.Items.Type, param.Items.Format)})
				continue
			}

			queryParam = append(queryParam, postman2.URLQueryParam{Key: param.Name, Value: buildQueryParamDefaultValue(param.Type, param.Format)})
		}
	}

	return queryParam

}

//buildQueryParamDefaultValue returns default values from a param type and format
func buildQueryParamDefaultValue(propType string, propFormat string) interface{} {

	if propType == "" {
		return ""
	}

	//Property has no example value : we set one by default
	if propType == "integer" {
		return 0
	}

	if propType == "string" {
		switch propFormat {
		case "date-time":
			return time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC).Format(time.RFC3339)
		default:
			return "string"
		}
	}

	return ""
}

//pathHasMethodWithTag checks if a swagger path if defined for a given method (GET, PUT, POST, etc.)
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
