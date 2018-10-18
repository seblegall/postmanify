package postman2_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Meetic/postmanify/postman2"
)

func TestNewCollection(t *testing.T) {

	collection := postman2.NewCollection("title", "desciption")

	assert.Equal(t, "title", collection.Info.Name)
	assert.Equal(t, "desciption", collection.Info.Description)
	assert.Equal(t, postman2.Schema, collection.Info.Schema)
}

func TestAddItem(t *testing.T) {

	collection := postman2.NewCollection("title", "desciption")

	collection.AddItem(postman2.APIItem{
		Name: "item request",
	}, "test")

	assert.Equal(t, "test", collection.Item[0].Name)
	assert.Equal(t, "item request", collection.Item[0].Item[0].Name)

	collection.AddItem(postman2.APIItem{
		Name: "item request 2",
	}, "test")

	assert.Equal(t, "test", collection.Item[0].Name)
	assert.Len(t, collection.Item, 1)
	assert.Len(t, collection.Item[0].Item, 2)
	assert.Equal(t, "item request 2", collection.Item[0].Item[1].Name)
}

func TestNewURL(t *testing.T) {

	dataset := []struct {
		input    string
		expected postman2.URL
	}{
		{
			input: "http://test.test.com/test/test",
			expected: postman2.URL{
				Raw:      "http://test.test.com/test/test",
				Protocol: "http",
				Host:     []string{"test", "test", "com"},
				Path:     []string{"test", "test"},
			},
		},
		{
			input: "https://{{test}}{{test}}{{test}}/test?test=true",
			expected: postman2.URL{
				Raw:      "https://{{test}}{{test}}{{test}}/test?test=true",
				Protocol: "https",
				Host:     []string{"{{test}}{{test}}{{test}}"},
				Path:     []string{"test?test=true"},
			},
		},
	}

	for _, data := range dataset {
		url := postman2.NewURL(data.input)
		assert.IsType(t, postman2.URL{}, url)
		assert.Equal(t, data.expected.Raw, url.Raw)
		assert.Equal(t, data.expected.Protocol, url.Protocol)
		assert.Equal(t, data.expected.Host, url.Host)
		assert.Equal(t, data.expected.Path, url.Path)

	}
}

func TestAddVariable(t *testing.T) {
	url := postman2.URL{
		Raw:      "http://test.test.com/test/test",
		Protocol: "http",
		Host:     []string{"test", "test", "com"},
		Path:     []string{"test", "test"},
	}

	url.AddVariable("test", "value")

	assert.Equal(t, "test", url.Variable[0].ID)
	assert.Equal(t, "value", url.Variable[0].Value)

	url.AddVariable("test2", "value2")

	assert.Equal(t, "test2", url.Variable[1].ID)
	assert.Equal(t, "value2", url.Variable[1].Value)
}
