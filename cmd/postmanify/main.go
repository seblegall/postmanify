package main

import (
	"flag"
	"io/ioutil"

	"github.com/Meetic/postmanify"
	"github.com/Meetic/postmanify/postman2"
)

var (
	swagSpecFilepath string
	pmanSpecFilepath string
	host string
)

func main() {

	flag.StringVar(&host, "host", "", `The hostname for the API`)
	flag.StringVar(&swagSpecFilepath, "f", "swagger.json", `The swagger file to convert`)
	flag.StringVar(&pmanSpecFilepath, "o", "postman_collection.json", `The postman collection file as output`)
	flag.Parse()

	conv := postmanify.NewConverter(postmanify.Config{
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

	if err := ioutil.WriteFile(pmanSpecFilepath, postman, 0644); err != nil {
		panic(err)
	}

}
