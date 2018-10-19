package swagger2_test

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Meetic/postmanify/swagger2"
)

const (
	swaggerFile = "fixtures/swagger_fixtures.json"
)

func TestNewSpecFromBytes(t *testing.T) {
	swagByte, _ := ioutil.ReadFile(swaggerFile)

	spec, err := swagger2.NewSpecificationFromBytes(swagByte)

	assert.Nil(t, err)
	assert.NotNil(t, spec)
	assert.Equal(t, "Swagger Petstore", spec.Info.Title)
	assert.Equal(t, "pet", spec.Paths["/pet"].Post.Tags[0])
}

func TestHasMethodWithTag(t *testing.T) {
	swagByte, _ := ioutil.ReadFile(swaggerFile)

	spec, _ := swagger2.NewSpecificationFromBytes(swagByte)

	path, _ := spec.Paths["/pet"]
	assert.True(t, path.HasMethodWithTag(http.MethodPost))
	assert.True(t, path.HasMethodWithTag(http.MethodPut))

	path2, _ := spec.Paths["/pet/findByStatus"]
	assert.True(t, path2.HasMethodWithTag(http.MethodGet))

	path3, _ := spec.Paths["/pet/{petId}"]
	assert.True(t, path3.HasMethodWithTag(http.MethodDelete))
}

func TestParseDefinition(t *testing.T) {
	dataset := []struct {
		input    swagger2.Schema
		expected string
	}{
		{
			input: swagger2.Schema{
				Ref: "#/definitions/User",
			},
			expected: "User",
		},
	}

	for _, data := range dataset {
		assert.Equal(t, data.expected, data.input.ParseDefinition())
	}

}
