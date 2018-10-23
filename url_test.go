package postmanify

import (
	"fmt"
	"testing"

	"github.com/Meetic/postmanify/postman2"
	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
)

func TestBuildPostmanURL(t *testing.T) {

	var enumInt []interface{}
	enumInt = append(enumInt, 1)
	enumInt = append(enumInt, 2)
	enumInt = append(enumInt, 3)

	var enumString []interface{}
	enumString = append(enumString, "test1")
	enumString = append(enumString, "test2")
	enumString = append(enumString, "test3")

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
				url: "/test/test/{test}",
				operation: &spec.Operation{
					OperationProps: spec.OperationProps{
						Parameters: []spec.Parameter{
							spec.Parameter{
								ParamProps: spec.ParamProps{
									Name: "test",
									In:   "query",
								},
								SimpleSchema: spec.SimpleSchema{
									Type:    "string",
									Example: "test",
								},
							},
							spec.Parameter{
								ParamProps: spec.ParamProps{
									Name: "testEnum",
									In:   "query",
								},
								SimpleSchema: spec.SimpleSchema{
									Type: "string",
								},
								CommonValidations: spec.CommonValidations{
									Enum: enumString,
								},
							},
							spec.Parameter{
								ParamProps: spec.ParamProps{
									Name: "test2",
									In:   "query",
								},
								SimpleSchema: spec.SimpleSchema{
									Type: "array",
									Items: &spec.Items{
										SimpleSchema: spec.SimpleSchema{
											Type: "integer",
										},
										CommonValidations: spec.CommonValidations{
											Enum: enumInt,
										},
									},
								},
							},
						},
					},
				},
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
				Query: []postman2.URLQueryParam{
					postman2.URLQueryParam{
						Key:   "test",
						Value: "test",
					},
					postman2.URLQueryParam{
						Key:   "testEnum",
						Value: "test1",
					},
					postman2.URLQueryParam{
						Key:   "test2",
						Value: 1,
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
		//

		for i, param := range data.expected.Query {
			assert.Equal(t, param.Key, url.Query[i].Key)
			assert.Equal(t, param.Value, url.Query[i].Value)
		}

	}

}

func TestBuildQueryParams(t *testing.T) {

	var enumInt []interface{}
	enumInt = append(enumInt, 1)
	enumInt = append(enumInt, 2)
	enumInt = append(enumInt, 3)

	var enumString []interface{}
	enumString = append(enumString, "test1")
	enumString = append(enumString, "test2")
	enumString = append(enumString, "test3")

	dataset := []struct {
		input    *spec.Operation
		expected []postman2.URLQueryParam
	}{
		{
			input: &spec.Operation{
				OperationProps: spec.OperationProps{
					Parameters: []spec.Parameter{
						spec.Parameter{
							ParamProps: spec.ParamProps{
								Name: "test",
								In:   "query",
							},
							SimpleSchema: spec.SimpleSchema{
								Type:    "string",
								Example: "test",
							},
						},
						spec.Parameter{
							ParamProps: spec.ParamProps{
								Name: "testEnum",
								In:   "query",
							},
							SimpleSchema: spec.SimpleSchema{
								Type: "string",
							},
							CommonValidations: spec.CommonValidations{
								Enum: enumString,
							},
						},
						spec.Parameter{
							ParamProps: spec.ParamProps{
								Name: "test2",
								In:   "query",
							},
							SimpleSchema: spec.SimpleSchema{
								Type: "array",
								Items: &spec.Items{
									SimpleSchema: spec.SimpleSchema{
										Type: "integer",
									},
									CommonValidations: spec.CommonValidations{
										Enum: enumInt,
									},
								},
							},
						},
					},
				},
			},
			expected: []postman2.URLQueryParam{
				postman2.URLQueryParam{
					Key:   "test",
					Value: "test",
				},
				postman2.URLQueryParam{
					Key:   "testEnum",
					Value: "test1",
				},
				postman2.URLQueryParam{
					Key:   "test2",
					Value: 1,
				},
			},
		},
	}

	for _, data := range dataset {

		fmt.Println(data.input.Parameters[1].Enum)
		assert.Equal(t, data.expected, buildQueryParams(data.input))
	}
}
