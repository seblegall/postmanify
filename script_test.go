package postmanify

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Meetic/postmanify/postman2"
	"github.com/go-openapi/spec"
)

func TestbuildPostmanScript(t *testing.T) {

	var script []interface{}
	script = append(script, "test")
	script = append(script, "test2")

	dataset := []struct {
		input    spec.Extensions
		expected postman2.Script
	}{
		{
			input: spec.Extensions{
				"x-postman-script": script,
				"x-test":           script,
			},
			expected: postman2.Script{
				Type: scriptType,
				Exec: []string{"test", "test2"},
			},
		},
	}

	for _, data := range dataset {
		script := buildPostmanScript(data.input)
		for i, e := range data.expected.Exec {
			assert.Equal(t, e, script.Exec[i])
		}
	}
}
