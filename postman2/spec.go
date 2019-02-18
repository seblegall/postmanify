package postman2

import (
	"regexp"
	"strings"
)

const (
	//Schema represents the schema version of postman collection file
	Schema = "https://schema.getpostman.com/json/collection/v2.0.0/collection.json"
)

//Collection represents a Postman Collection
type Collection struct {
	Info CollectionInfo `json:"info"`
	Item []FolderItem   `json:"item"`
}

//NewCollection creates a Postman Collection using the title and the description
func NewCollection(title, description string) Collection {
	return Collection{
		Info: CollectionInfo{
			Name:        title,
			Description: description,
			Schema:      Schema,
		},
	}
}

//AddItem add a Postman item to the collection
func (col *Collection) AddItem(newItem APIItem, folder string) {
	for i, itemFolder := range col.Item {
		if itemFolder.Name == folder {
			col.Item[i].Item = append(itemFolder.Item, newItem)
			return
		}
	}

	newFolder := FolderItem{
		Name: folder,
		Item: []APIItem{newItem},
	}
	col.Item = append(col.Item, newFolder)
}

//CollectionInfo represents a Collection description
type CollectionInfo struct {
	Name        string `json:"name,omitempty"`
	PostmanID   string `json:"_postman_id,omitempty"`
	Description string `json:"description,omitempty"`
	Schema      string `json:"schema,omitempty"`
}

//FolderItem represents a Postman folder part of a collection
type FolderItem struct {
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Item        []APIItem `json:"item,omitempty"`
}

//APIItem represents a Postman request
type APIItem struct {
	Name    string  `json:"name,omitempty"`
	Event   []Event `json:"event,omitempty"`
	Request Request `json:"request,omitempty"`
}

//Event represents a Postman Event aka a post-run script
type Event struct {
	Listen string `json:"listen"`
	Script Script `json:"script"`
}

//Script represents a Postman Script
type Script struct {
	Type string   `json:"type,omitempty"`
	Exec []string `json:"exec,omitmpety"`
}

//Request represents a Postman Request
type Request struct {
	URL         URL         `json:"url,omitempty"`
	Method      string      `json:"method,omitempty"`
	Header      []Header    `json:"header,omitempty"`
	Body        RequestBody `json:"body,omitempty"`
	Description string      `json:"description,omitempty"`
}

//URL represents a Postman URL, part from the request
type URL struct {
	Raw      string            `json:"raw,omitempty"`
	Protocol string            `json:"protocol,omitempty"`
	Auth     map[string]string `json:"auth"`
	Host     []string          `json:"host,omitempty"`
	Path     []string          `json:"path,omitempty"`
	Variable []URLVariable     `json:"variable,omitempty"`
	Query    []URLQueryParam   `json:"query,omitempty"`
}

//URL represents a Postman URL variable, part from the URL
type URLVariable struct {
	Value interface{} `json:"value,omitempty"`
	ID    string      `json:"id,omitempty"`
}

//NewURL creates a Postman URL from a rowURL
//It extracts the protocol, the host and the path
func NewURL(rawURL string) URL {
	rawURL = strings.TrimSpace(rawURL)
	url := URL{Raw: rawURL, Variable: []URLVariable{}}
	rx := regexp.MustCompile(`^([a-z]+)://([^/]+)/(.*)$`)
	rs := rx.FindAllStringSubmatch(rawURL, -1)

	if len(rs) > 0 {
		for _, m := range rs {
			url.Protocol = m[1]

			hostname := m[2]
			hostnameParts := strings.Split(hostname, ".")
			url.Host = hostnameParts

			path := m[3]
			pathParts := strings.Split(path, "/")
			url.Path = pathParts
		}
	}

	return url
}

//AddVariable add a Postman URL variable to the URL
func (url *URL) AddVariable(key string, value interface{}) {
	variable := URLVariable{ID: key, Value: value}
	url.Variable = append(url.Variable, variable)
}

//AddQueryParam add a query param to the URL
func (url *URL) AddQueryParam(param URLQueryParam) {
	url.Query = append(url.Query, param)
}

//URLQueryParam represents a Postman URL query param
type URLQueryParam struct {
	Key   string      `json:"key,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

//Header represents a header, part from the Postman request
type Header struct {
	Key         string `json:"key,omitempty"`
	Value       string `json:"value,omitempty"`
	Description string `json:"description,omitempty"`
}

//RequestBody represents the Postman request's body
type RequestBody struct {
	Mode       string            `json:"mode,omitempty"`
	URLEncoded []URLEncodedParam `json:"urlencoded,omitempty"`
	FormData   []FormData        `json:"formdata,omitempty"`
	Raw        string            `json:"raw,omitempty"`
}

//FormData represents the request body formatted as formdata
type FormData struct {
	Key     string `json:"key,omitempty"`
	Value   string `json:"value,omitempty"`
	Type    string `json:"type,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
}

//URLEncodedParam represents the request body formatted as URl encoded
type URLEncodedParam struct {
	Key     string `json:"key,omitempty"`
	Value   string `json:"value,omitempty"`
	Type    string `json:"type,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
}
