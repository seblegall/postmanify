# Postmanify

Postmanify is a simple tool allowing you to convert Swagger \(renamed Open-API\) spec files into Postman collections.

Postmanify is the only swagger to postman converter able to create Postman POST/PUT request with a pre-filled json body.

## Installation

Postmanify is available for all OS. You only have to execute this command to download the right binary for your system.

```text
curl -sf https://raw.githubusercontent.com/seblegall/postmanify/master/install.sh | sh
```

## Usage

```sh
$ postmanify --help
Usage of postmanify:
  -f string
        The swagger file to convert (default "swagger.json")
  -host string
        The hostname for the API
  -o string
        The postman collection file as output (default "postman_collection.json")
```

## Features

### Postman variables

By using the `{your_variable}` notation in your swagger file, Postmanify is able to create postman environnement variable automatically in the generated postman collection.

A good use of variables may be in the path definitions :

```json
"paths": {
    "/my/app/route/{param}": {
    },
```

By doing so, each people importing the generated postman collection will be able to customize its own value for the environnement variable called `param` directly from postman.

### Postman scripts

With Postmanify you can also document Postman script directly in your swagger file by using the special key `x-postman-script`. This way, each people importing the generated postman collection will benefits  of the postman scripts already configured.

A good usage of Postman scripts may be done to populate environnement variable with value received from a previous call.

The example bellow populate environnement variable called `token` with the content of the key `access_token` from the return payload.

```json
{
    "paths": {
        "/oauth/accesstokens": {
            "post": {
                "x-postman-script": [
                "var jsonData = JSON.parse(responseBody);",
                "postman.setEnvironmentVariable(\"token\", jsonData.access_token);"
                ],
                "summary": "Create a new session and return an accesstoken",
            }
        },
    }
}
```

Once the env var populated, it's easy to reuse it by using the postman variable notation in the value of a field in a swagger spec, such as an Authorization header.

### Auto-generated request body

Postmanify is the only swagger to postman converter able to generate request body directly from the swagger specs.

To do so, Postmanify rebuild the json request body from the body `parameters` specified in the swagger specs. It attributes a the best value to each parameter based on :

* The required field
* The example field
* The enum field
* The type

In other words, Postmanify first check if the field is required or not. It only build payloads with required parameters.

Then, it checks if the `example` is filled and takes the value of this field as a default value.

Else, if the `enum` field is filled, it takes the first value as a default value.

Else, it generates a default value based on the field type. For example, for a string field, it will used "string" as a default value.

Bellow is an example of Swagger request body converted as a json payload. In this example the POST request is made to create a new user on a given API.

Swagger spec : 

```json
"parameters": [
    {
    "description": "User object",
    "name": "body",
    "in": "body",
    "required": true,
    "schema": {
        "type": "object",
        "required": [
            "user",
            "metas"
        ],
        "properties": {
            "members": {
                "type": "object",
                "properties": {
                    "profile": {
                        "type": "object",
                        "properties": {
                            "gender" : {
                                "type": "string",
                                "enum": [
                                "F",
                                "M"
                                ]
                            },
                            "birth_date": {
                                "type": "string",
                                "format": "date-time",
                                "example": "1994-03-03T00:00:00+0100"
                            },
                            "nickname" : {
                                "type": "string",
                            },
                            "email": {
                                "type": "string",
                                "example": "test10@example.com"
                            },
                            "password": {
                                "type": "string",
                                "example": "123456aA"
                            }
                        }
                    },
                }
            },
            "metas": {
                "type": "object",
                "properties": {
                    "marketing_code": {
                        "type": "string",
                        "example": "000001"
                    },
                }
            }
        }
    }
    }
]
```

Generated Json Payload :

```json
{
    "members": {
        "profile": {
            "gender": "F",
            "birth_date": "1994-03-03T00:00:00+0100",
            "nickname": "string",
            "email": "test@example.com",
            "password": "123456aA"
        },
    },
    "metas": {
        "marketing_code": "000001",
    }
}
```

* `gender` takes the first value of the enum : "F"
* `birth_date` takes the value in example
* `nickname` falls back on the default value "string"



## Developers guide

Postmanify may also be used as a go package. Bellow is an example of how to integrate Postmanify as an external dependency.

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