# Postmanify

A simple Go library helping you to convert Swagger 2 spec document into Postman 2 collection

## Warning

This lib is under development. Use it at your own risk.

## Usage

```go
package main

import (
	"io/ioutil"

	"github.com/Meetic/postmanify"
	"github.com/Meetic/postmanify/postman2"
)

const (
	swagSpecFilepath = "swagger.json"
	pmanSpecFilepath = "postman.json"
	hostPrefix = "prefix"
	hostSuffix = "suffix"
)

func main() {
	conv := postmanify.NewConverter(postmanify.Config{
    		HostnamePrefix: hostPrefix,
    		HostnameSuffix: hostSuffix,
    		PostmanHeaders: map[string]postman2.Header{
    			"Authorization": {
    				Key:   "Authorization",
    				Value: "Bearer {{my_access_token}}"},
    		},
    	})
    
    	swag, _ := ioutil.ReadFile(swagSpecFilepath)
    	postman, err := conv.Convert(swag)
    	if err != nil {
    		panic(err)
    	}
    
    	ioutil.WriteFile(pmanSpecFilepath, postman, 0644)

}
```