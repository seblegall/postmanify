package postmanify

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
)

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
		{
			input: map[string]spec.Schema{
				"id": spec.Schema{
					SwaggerSchemaProps: spec.SwaggerSchemaProps{
						Example: 1234,
					},
				},
				"result": spec.Schema{
					SchemaProps: spec.SchemaProps{
						Type: spec.StringOrArray{"array"},
						Items: &spec.SchemaOrArray{
							Schema: &spec.Schema{
								SchemaProps: spec.SchemaProps{
									Type: spec.StringOrArray{"string"},
								},
							},
						},
					},
				},
			},
			expected: indentJSON(`{"id":1234,"result":["string"]}`),
		},
		{
			input: map[string]spec.Schema{
				"id": spec.Schema{
					SwaggerSchemaProps: spec.SwaggerSchemaProps{
						Example: 1234,
					},
				},
				"result": spec.Schema{
					SchemaProps: spec.SchemaProps{
						Type: spec.StringOrArray{"array"},
						Items: &spec.SchemaOrArray{
							Schema: &spec.Schema{
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
						},
					},
				},
			},
			expected: indentJSON(`{"id":1234,"result":[{"test": 1234}]}`),
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
