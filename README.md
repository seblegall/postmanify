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
