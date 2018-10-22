package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"

	"github.com/Meetic/postmanify"
	"github.com/Meetic/postmanify/postman2"
)

const (
	swagSpecFilepath         = "swagger.json"
	swagSpecExtantedFilepath = "swagger_extended.json"
	pmanSpecFilepath         = "postman.json"
)

func main() {
	conv := postmanify.NewConverter(postmanify.Config{
		HostnamePrefix: "prefix.",
		HostnameSuffix: ".suffix.com",
		PostmanHeaders: []postman2.Header{{
			Key:   "Authorization",
			Value: "Bearer {{my_access_token}}"}},
	})

	swag, _ := ioutil.ReadFile(swagSpecFilepath)
	postman, err := conv.Convert(swag)
	if err != nil {
		panic(err)
	}

	ioutil.WriteFile(pmanSpecFilepath, postman, 0644)

	writeExtanded()

}

func writeExtanded() {
	specDoc, err := loads.Spec(swagSpecFilepath)

	if err != nil {
		panic(err)
	}

	specDocExpand, err := specDoc.Expanded(&spec.ExpandOptions{
		SkipSchemas:         false,
		ContinueOnError:     true,
		AbsoluteCircularRef: true,
	})
	if err != nil {
		panic(err)
	}

	postman, _ := json.MarshalIndent(specDocExpand.Spec(), "", "  ")

	ioutil.WriteFile(swagSpecExtantedFilepath, postman, 0644)

}
