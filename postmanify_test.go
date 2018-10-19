package postmanify

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/Meetic/postmanify/postman2"
	"github.com/Meetic/postmanify/swagger2"

	"github.com/stretchr/testify/assert"
)

func TestNewConverter(t *testing.T) {
	conv := NewConverter(Config{})

	assert.NotNil(t, conv)
}

func TestBuildPostmanURL(t *testing.T) {
	dataset := []struct {
		input struct {
			cfg      Config
			url      string
			endpoint swagger2.Endpoint
		}
		expected postman2.URL
	}{
		{
			input: struct {
				cfg      Config
				url      string
				endpoint swagger2.Endpoint
			}{
				cfg: Config{
					Hostname:       "hostname",
					HostnamePrefix: "prefix.",
					HostnameSuffix: ".suffix.com",
					BasePath:       "/test/",
					Schema:         "http",
				},
				url:      "/test / test /  {test}",
				endpoint: swagger2.Endpoint{},
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
		url := conv.buildPostmanURL(data.input.url, data.input.endpoint)

		assert.Equal(t, data.expected.Raw, url.Raw)
		assert.Equal(t, data.expected.Protocol, url.Protocol)
		assert.Equal(t, data.expected.Host, url.Host)
		assert.Equal(t, data.expected.Variable[0].ID, url.Variable[0].ID)

	}

}

func TestBuildProperties(t *testing.T) {

	definitions := map[string]swagger2.Definition{
		"Test": swagger2.Definition{
			Type: "object",
			Properties: map[string]swagger2.Property{
				"test": swagger2.Property{
					Example: 1234,
				},
			},
		},
	}

	dataset := []struct {
		input    map[string]swagger2.Property
		expected string
	}{
		{
			input: map[string]swagger2.Property{
				"id": swagger2.Property{
					Example: 1234,
				},
				"username": swagger2.Property{
					Example: "john",
				},
			},
			expected: indentJSON(`{"id":1234,"username":"john"}`),
		},
		{
			input: map[string]swagger2.Property{
				"createdAt": swagger2.Property{
					Type:   "string",
					Format: "date-time",
				},
				"id": swagger2.Property{
					Example: 1234,
				},
				"username": swagger2.Property{
					Type: "string",
				},
			},
			expected: indentJSON(`{"createdAt":"1994-03-03T00:00:00+0100","id":1234,"username":"string"}`),
		},
		{
			input: map[string]swagger2.Property{
				"createdAt": swagger2.Property{
					Type:   "string",
					Format: "date-time",
				},
				"id": swagger2.Property{
					Example: 1234,
				},
				"test": swagger2.Property{
					Ref: "#/definitions/Test",
				},
				"username": swagger2.Property{
					Type: "string",
				},
			},
			expected: indentJSON(`{"createdAt":"1994-03-03T00:00:00+0100","id":1234, "test":{"test":1234},"username":"string"}`),
		},
	}

	for _, data := range dataset {
		conv := NewConverter(Config{})
		conv.definitions = definitions
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
