package postmanify

import (
	"encoding/json"
	"strings"

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
