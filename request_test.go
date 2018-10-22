package postmanify

import (
	"testing"

	"github.com/Meetic/postmanify/postman2"
	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
)

func TestBuildPostmanBody(t *testing.T) {
	dataset := []struct {
		input    *spec.Operation
		expected postman2.RequestBody
	}{
		{
			input: &spec.Operation{
				OperationProps: spec.OperationProps{
					Parameters: []spec.Parameter{
						spec.Parameter{
							ParamProps: spec.ParamProps{
								In:       "body",
								Required: true,
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type: spec.StringOrArray{"object"},
										Properties: map[string]spec.Schema{
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
									},
								},
							},
						},
					},
				},
			},
			expected: postman2.RequestBody{
				Mode: "raw",
				Raw:  indentJSON(`{"id":1234,"username":"john"}`),
			},
		},
		{
			input: &spec.Operation{
				OperationProps: spec.OperationProps{
					Parameters: []spec.Parameter{
						spec.Parameter{
							ParamProps: spec.ParamProps{
								In:       "body",
								Required: true,
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type: spec.StringOrArray{"object"},
										Properties: map[string]spec.Schema{
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
									},
								},
							},
						},
					},
				},
			},
			expected: postman2.RequestBody{
				Mode: "raw",
				Raw:  indentJSON(`{"id":1234,"result":["string"]}`),
			},
		},
	}

	for _, data := range dataset {

		conv := NewConverter(Config{})

		requestBody := conv.buildPostmanBody(data.input)

		assert.Equal(t, data.expected.Mode, requestBody.Mode)
		assert.Equal(t, data.expected.Raw, requestBody.Raw)
	}

}
