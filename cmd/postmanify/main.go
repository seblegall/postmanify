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
	hostPrefix       string
	hostSuffix       string
)

func main() {

	flag.StringVar(&hostPrefix, "host-prefix", "", `A prefix to put before the globale hostname`)
	flag.StringVar(&hostSuffix, "host-suffix", "", `A suffix to put after the globale hostname`)
	flag.StringVar(&swagSpecFilepath, "f", "swagger.json", `The swagger file to convert`)
	flag.StringVar(&pmanSpecFilepath, "o", "postman_collection.json", `The postman collection file as output`)
	flag.Parse()

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
