package postmanify

import (
	"strings"

	"github.com/seblegall/postmanify/postman2"
	"github.com/go-openapi/spec"
)

const (
	scriptType = "text/javascript"
)

//buildPostmanScript creates a postman js script from a "x-postman-script" swagger extension
func buildPostmanScript(extensions spec.Extensions) postman2.Script {

	if s, ok := extensions.GetString("x-postman-script"); ok {
		return postman2.Script{
			Type: scriptType,
			Exec: strings.Split(s, "\n"),
		}
	}

	if s, ok := extensions.GetStringSlice("x-postman-script"); ok {
		return postman2.Script{
			Type: scriptType,
			Exec: s,
		}
	}

	return postman2.Script{}
}
