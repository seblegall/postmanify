package postmanify

import (
	"fmt"
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

		fmt.Println(conv.Config)

		assert.Equal(t, data.expected.Raw, url.Raw)
		assert.Equal(t, data.expected.Protocol, url.Protocol)
		assert.Equal(t, data.expected.Host, url.Host)
		assert.Equal(t, data.expected.Variable[0].ID, url.Variable[0].ID)

	}

}
