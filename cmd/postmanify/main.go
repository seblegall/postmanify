package main

import (
	"io/ioutil"

	"github.com/Meetic/postmanify"
	"github.com/Meetic/postmanify/postman2"
)

const (
	swagSpecFilepath = "swagger.json"
	pmanSpecFilepath = "postman.json"
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

}
