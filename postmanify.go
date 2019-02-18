package postmanify

import (
	"encoding/json"
	"strings"

	"github.com/seblegall/postmanify/postman2"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
)

//Config represents an converter configuration
type Config struct {
	//Hostname may be used if you want to override the swagger defined hostname
	Hostname       string
	//Schema represents the protocol. It may be used if you want to override the one defined in the swagger file.
	Schema         string
	//BasePath may be used to override the swagger defined basPath.
	BasePath       string
	//PostmanHeaders represents a collection of header to add on each documented path before generating the corresponding postman collection.
	PostmanHeaders map[string]postman2.Header
}

//Converter represent a Swagger2.0 documentation to Postman 2.1 collections converter
type Converter struct {
	config Config
}


//NewConverter creates a new converter
func NewConverter(cfg Config) *Converter {
	return &Converter{
		config: cfg,
	}
}

//Convert converts a swagger specification to a postman collection.
//Convert expected a json input defined as a slice of byte, and returns a json, defined as a slice of byte
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
