package postmanify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"

	"github.com/Meetic/postmanify/postman2"
	"github.com/Meetic/postmanify/swagger2"
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
	Config Config
}

func NewConverter(cfg Config) *Converter {
	return &Converter{
		Config: cfg,
	}
}

func (c *Converter) Convert(swaggerSpec []byte) ([]byte, error) {

	swag, err := swagger2.NewSpecificationFromBytes(swaggerSpec)
	if err != nil {
		return nil, err
	}

	if c.Config.Hostname == "" {
		c.Config.Hostname = strings.TrimSpace(swag.Host)
	}

	if c.Config.BasePath == "" {
		c.Config.BasePath = strings.TrimSpace(swag.BasePath)
	}

	//if schema is not defined in config, we take the first one declared on the swagger specs.
	if c.Config.Schema == "" && len(swag.Schemes) >= 1 {
		c.Config.Schema = strings.TrimSpace(swag.Schemes[0])
	} else {
		c.Config.Schema = "http"
	}

	pman := postman2.NewCollection(strings.TrimSpace(swag.Info.Title), strings.TrimSpace(swag.Info.Description))

	if err := c.addUrls(swag.Paths, &pman); err != nil {
		return nil, err
	}

	return json.MarshalIndent(pman, "", "  ")

}

func (c *Converter) addUrls(paths map[string]swagger2.Path, pman *postman2.Collection) error {
	urls := []string{}
	for url := range paths {
		urls = append(urls, url)
	}
	sort.Strings(urls)

	for _, url := range urls {
		path := paths[url]

		if path.HasMethodWithTag(http.MethodGet) {
			pman.AddItem(c.buildPostmanItem(url, http.MethodGet, path.Get), strings.TrimSpace(path.Get.Tags[0]))
		}
		if path.HasMethodWithTag(http.MethodPatch) {
			pman.AddItem(c.buildPostmanItem(url, http.MethodPatch, path.Patch), strings.TrimSpace(path.Patch.Tags[0]))
		}
		if path.HasMethodWithTag(http.MethodPost) {
			pman.AddItem(c.buildPostmanItem(url, http.MethodPost, path.Post), strings.TrimSpace(path.Post.Tags[0]))
		}
		if path.HasMethodWithTag(http.MethodPut) {
			pman.AddItem(c.buildPostmanItem(url, http.MethodPut, path.Put), strings.TrimSpace(path.Put.Tags[0]))
		}
		if path.HasMethodWithTag(http.MethodDelete) {
			pman.AddItem(c.buildPostmanItem(url, http.MethodDelete, path.Delete), strings.TrimSpace(path.Delete.Tags[0]))
		}
	}

	return nil
}

func (c *Converter) buildPostmanItem(url, method string, endpoint swagger2.Endpoint) postman2.APIItem {

	return postman2.APIItem{
		Name:    url,
		Request: c.buildPostmanRequest(url, method, endpoint),
	}

}

func (c *Converter) buildPostmanRequest(url, method string, endpoint swagger2.Endpoint) postman2.Request {

	return postman2.Request{
		Method: strings.ToUpper(method),
		URL:    c.buildPostmanURL(url, endpoint),
		Header: c.buildPostmanHeaders(endpoint),
	}
}

func (c *Converter) buildPostmanHeaders(endpoint swagger2.Endpoint) []postman2.Header {
	headers := []postman2.Header{}
	if len(endpoint.Consumes) > 0 {
		if len(strings.TrimSpace(endpoint.Consumes[0])) > 0 {
			headers = append(headers, postman2.Header{
				Key:   "Content-Type",
				Value: strings.TrimSpace(endpoint.Consumes[0])})
		}
	}
	if len(endpoint.Produces) > 0 {
		if len(strings.TrimSpace(endpoint.Produces[0])) > 0 {
			headers = append(headers, postman2.Header{
				Key:   "Accept",
				Value: strings.TrimSpace(endpoint.Produces[0])})
		}
	}
	headers = append(headers, c.Config.PostmanHeaders...)

	return headers
}

func (c *Converter) buildPostmanURL(url string, endpoint swagger2.Endpoint) postman2.URL {

	//create hostname
	host := strings.TrimSpace(strings.Join([]string{
		c.Config.HostnamePrefix,
		c.Config.Hostname,
		c.Config.HostnameSuffix,
	}, ""))

	//create URI
	rawPostmanURL := strings.TrimSpace(strings.Join([]string{
		host,
		c.Config.BasePath,
		strings.TrimSpace(url),
	}, "/"))
	rx1 := regexp.MustCompile(`/+`)
	rawPostmanURL = rx1.ReplaceAllString(rawPostmanURL, "/")
	rx2 := regexp.MustCompile(`^/+`)
	rawPostmanURL = rx2.ReplaceAllString(rawPostmanURL, "")

	//Add schema
	rawPostmanURL = strings.Join([]string{c.Config.Schema, rawPostmanURL}, "://")

	//Replace URI parameters
	fmt.Printf("tmp postman url : %s \n", rawPostmanURL)
	rx3 := regexp.MustCompile(`(^|[^\{])\{([^\/\{\}]+)\}([^\}]|$)`)
	rawPostmanURL = rx3.ReplaceAllString(rawPostmanURL, "$1{{$2}}$3")
	fmt.Printf("tmp postman url (after regex) : %s \n", rawPostmanURL)

	postmanURL := postman2.NewURL(rawPostmanURL)

	// Set Default URL Path Parameters
	rx4 := regexp.MustCompile(`^\s*({{(.+)}})\s*$`)

	for _, part := range postmanURL.Path {
		rs4 := rx4.FindAllStringSubmatch(part, -1)
		fmt.Printf("test url param : %s \n", rs4)
		if len(rs4) > 0 {
			baseVariable := rs4[0][2]
			var defaultValue interface{}
			for _, parameter := range endpoint.Parameters {
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
