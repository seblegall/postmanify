package postmanify

import (
	"net/http"
	"regexp"
	"sort"
	"strings"

	"github.com/Meetic/postmanify/postman2"
	"github.com/go-openapi/spec"
)

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
