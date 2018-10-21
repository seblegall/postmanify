package postmanify

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/Meetic/postmanify/postman2"
	"github.com/go-openapi/spec"

	"github.com/stretchr/testify/assert"
)

func TestNewConverter(t *testing.T) {
	conv := NewConverter(Config{})

	assert.NotNil(t, conv)
}

func TestBuildPostmanURL(t *testing.T) {
	dataset := []struct {
		input struct {
			cfg       Config
			url       string
			operation *spec.Operation
		}
		expected postman2.URL
	}{
		{
			input: struct {
				cfg       Config
				url       string
				operation *spec.Operation
			}{
				cfg: Config{
					Hostname:       "hostname",
					HostnamePrefix: "prefix.",
					HostnameSuffix: ".suffix.com",
					BasePath:       "/test/",
					Schema:         "http",
				},
				url:       "/test / test /  {test}",
				operation: &spec.Operation{},
			},
			expected: postman2.URL{
				Raw:      "http://prefix.hostname.suffix.com/test/test/test/{{test}}",
				Protocol: "http",
				Host:     []string{"prefix", "hostname", "suffix", "com"},
				Path:     []string{"test", "test", "test", "{{test}}"},
				Variable: []postman2.URLVariable{
					{
						ID: "test",
					},
				},
			},
		},
	}

	for _, data := range dataset {
		conv := NewConverter(data.input.cfg)
		url := conv.buildPostmanURL(data.input.url, data.input.operation)

		assert.Equal(t, data.expected.Raw, url.Raw)
		assert.Equal(t, data.expected.Protocol, url.Protocol)
		assert.Equal(t, data.expected.Host, url.Host)
		assert.Equal(t, data.expected.Variable[0].ID, url.Variable[0].ID)

	}

}

func TestBuildProperties(t *testing.T) {

	dataSlice := []string{"ok", "nok"}
	var interfaceSlice []interface{} = make([]interface{}, len(dataSlice))
	for i, d := range dataSlice {
		interfaceSlice[i] = d
	}

	dataset := []struct {
		input    map[string]spec.Schema
		expected string
	}{
		{
			input: map[string]spec.Schema{
				"id": spec.Schema{
					SwaggerSchemaProps: spec.SwaggerSchemaProps{
						Example: 1234,
					},
				},
				"username": spec.Schema{
					SwaggerSchemaProps: spec.SwaggerSchemaProps{
						Example: "john",
					},
				},
			},
			expected: indentJSON(`{"id":1234,"username":"john"}`),
		},
		{
			input: map[string]spec.Schema{
				"createdAt": spec.Schema{
					SchemaProps: spec.SchemaProps{
						Type:   spec.StringOrArray{"string"},
						Format: "date-time",
					},
				},
				"id": spec.Schema{
					SwaggerSchemaProps: spec.SwaggerSchemaProps{
						Example: 1234,
					},
				},
				"username": spec.Schema{
					SchemaProps: spec.SchemaProps{
						Type: spec.StringOrArray{"string"},
					},
				},
			},
			expected: indentJSON(`{"createdAt":"2009-11-17T20:34:58Z","id":1234,"username":"string"}`),
		},
		{
			input: map[string]spec.Schema{
				"createdAt": spec.Schema{
					SchemaProps: spec.SchemaProps{
						Type:   spec.StringOrArray{"string"},
						Format: "date-time",
					},
				},
				"id": spec.Schema{
					SwaggerSchemaProps: spec.SwaggerSchemaProps{
						Example: 1234,
					},
				},
				"test": spec.Schema{
					SchemaProps: spec.SchemaProps{
						Type: spec.StringOrArray{"object"},
						Properties: map[string]spec.Schema{
							"test": spec.Schema{
								SchemaProps: spec.SchemaProps{
									Type: spec.StringOrArray{"integer"},
								},
								SwaggerSchemaProps: spec.SwaggerSchemaProps{
									Example: 1234,
								},
							},
						},
					},
				},
				"username": spec.Schema{
					SchemaProps: spec.SchemaProps{
						Type: spec.StringOrArray{"string"},
					},
				},
			},
			expected: indentJSON(`{"createdAt":"2009-11-17T20:34:58Z","id":1234, "test":{"test":1234},"username":"string"}`),
		},
		{
			input: map[string]spec.Schema{
				"id": spec.Schema{
					SwaggerSchemaProps: spec.SwaggerSchemaProps{
						Example: 1234,
					},
				},
				"status": spec.Schema{
					SchemaProps: spec.SchemaProps{
						Type: spec.StringOrArray{"string"},
						Enum: interfaceSlice,
					},
				},
			},
			expected: indentJSON(`{"id":1234,"status":"ok"}`),
		},
	}

	for _, data := range dataset {
		conv := NewConverter(Config{})
		assert.Equal(t, data.expected, string(conv.buildProperties(data.input)))

	}
}

func indentJSON(s string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(s), "", "\t")
	if err != nil {
		panic(err)
	}
	return string(out.Bytes())
}
