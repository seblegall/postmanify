package postmanify

import (
	"strings"

	"github.com/Meetic/postmanify/postman2"
	"github.com/go-openapi/spec"
)

const (
	scriptType = "text/javascript"
)

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
