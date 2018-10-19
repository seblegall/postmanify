package swagger2

import (
	"encoding/json"
	"regexp"
	"strings"
)

// Specification represents a Swagger 2.0 specification.
type Specification struct {
	Swagger     string                `json:"swagger,omitempty"`
	Host        string                `json:"host,omitempty"`
	Info        *Info                 `json:"info,omitempty"`
	BasePath    string                `json:"basePath,omitempty"`
	Schemes     []string              `json:"schemes,omitempty"`
	Paths       map[string]Path       `json:"paths,omitempty"`
	Definitions map[string]Definition `json:"definitions,omitempty"`
}

// NewSpecificationFromBytes returns a Swagger Specification from a byte array.
func NewSpecificationFromBytes(data []byte) (Specification, error) {
	spec := Specification{}
	err := json.Unmarshal(data, &spec)
	return spec, err
}

// Info represents a Swagger 2.0 spec info object.
type Info struct {
	Description    string `json:"description,omitempty"`
	Version        string `json:"version,omitempty"`
	Title          string `json:"title,omitempty"`
	TermsOfService string `json:"termsOfService,omitempty"`
}

// Path represents a Swagger 2.0 spec path object.
type Path struct {
	Get    Endpoint `json:"get,omitempty"`
	Patch  Endpoint `json:"patch,omitempty"`
	Post   Endpoint `json:"post,omitempty"`
	Put    Endpoint `json:"put,omitempty"`
	Delete Endpoint `json:"delete,omitempty"`
	Ref    string   `json:"$ref,omitempty"`
}

func (p *Path) HasMethodWithTag(method string) bool {
	method = strings.TrimSpace(strings.ToLower(method))
	switch method {
	case "get":
		if &p.Get != nil && len(p.Get.Tags) > 0 && len(strings.TrimSpace(p.Get.Tags[0])) > 0 {
			return true
		}
	case "patch":
		if &p.Patch != nil && len(p.Patch.Tags) > 0 && len(strings.TrimSpace(p.Patch.Tags[0])) > 0 {
			return true
		}
	case "post":
		if &p.Post != nil && len(p.Post.Tags) > 0 && len(strings.TrimSpace(p.Post.Tags[0])) > 0 {
			return true
		}
	case "put":
		if &p.Put != nil && len(p.Put.Tags) > 0 && len(strings.TrimSpace(p.Put.Tags[0])) > 0 {
			return true
		}
	case "delete":
		if &p.Delete != nil && len(p.Delete.Tags) > 0 && len(strings.TrimSpace(p.Delete.Tags[0])) > 0 {
			return true
		}
	}
	return false
}

// Endpoint represents a Swagger 2.0 spec endpoint object.
type Endpoint struct {
	Tags        []string            `json:"tags,omitempty"`
	Summary     string              `json:"summary,omitempty"`
	OperationID string              `json:"operationId,omitempty"`
	Description string              `json:"description,omitempty"`
	Consumes    []string            `json:"consumes,omitempty"`
	Produces    []string            `json:"produces,omitempty"`
	Parameters  []Parameter         `json:"parameters,omitempty"`
	Responses   map[string]Response `json:"responses,omitempty"`
}

type Response struct {
	Description string            `json:"description,omitempty"`
	Schema      Schema            `json:"schema,omitempty"`
	Headers     map[string]Header `json:"headers,omitempty"`
}

type Schema struct {
	Type string `json:"type,omitempty"`
	Ref  string `json:"$ref,omitempty"`
}

func ParseDefinition(s string) string {
	// "#/definitions/User"
	rx := regexp.MustCompile(`^\#\/definitions\/((.+))$`)
	result := rx.FindStringSubmatch(s)
	if len(result) > 2 {
		return result[1]
	}

	return ""
}

type Header struct {
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
}

type Definition struct {
	Type       string              `json:"type,omitempty"`
	Properties map[string]Property `json:"properties,omitempty"`
}

type Property struct {
	Description string      `json:"description,omitempty"`
	Format      string      `json:"format,omitempty"`
	Items       Items       `json:"items,omitempty"`
	Type        string      `json:"type,omitempty"`
	Ref         string      `json:"$ref,omitempty"`
	Example     interface{} `json:"example,omitempty"`
}

type Items struct {
	Type string `json:"type,omitempty"`
	Ref  string `json:"$ref,omitempty"`
}

// Parameter represents a Swagger 2.0 spec parameter object.
type Parameter struct {
	Name        string      `json:"name,omitempty"`
	Type        string      `json:"type,omitempty"`
	In          string      `json:"in,omitempty"`
	Description string      `json:"description,omitempty"`
	Required    bool        `json:"required,omitempty"`
	Default     interface{} `json:"default,omitempty"`
	Schema      Schema      `json:"schema,omitempty"`
}
