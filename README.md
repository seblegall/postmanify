# Postmanify

A simple Go library helping you to convert Swagger 2 spec document into Postman 2 collection.

Postmanify also comes with a binary ready to be used.

## Warning

Postmanify is still under development. First release is an alpha version.

## Lib usage

```go
package main

import (
	"io/ioutil"

	"github.com/seblegall/postmanify"
	"github.com/seblegall/postmanify/postman2"
)

const (
	swagSpecFilepath = "swagger.json"
	pmanSpecFilepath = "postman.json"
	hostname = "my.api.com"
)

func main() {
	conv := postmanify.NewConverter(postmanify.Config{
    		Hostname: hostname,
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

## Binary usage

### install

```sh
curl -sf https://raw.githubusercontent.com/seblegall/postmanify/master/install.sh | sh
```

### usage 

```sh
  -f string
        The swagger file to convert (default "swagger.json")
  -host string
        The hostname for the API
  -o string
        The postman collection file as output (default "postman_collection.json")
```