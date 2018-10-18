package postman2

import (
	"regexp"
	"strings"
)

const (
	//Schema represents the schema version of postman collection file
	Schema = "https://schema.getpostman.com/json/collection/v2.0.0/collection.json"
)

type Collection struct {
	Info CollectionInfo `json:"info"`
	Item []FolderItem   `json:"item"`
}

func NewCollection(title, description string) Collection {
	return Collection{
		Info: CollectionInfo{
			Name:        title,
			Description: description,
			Schema:      Schema,
		},
	}
}

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

type CollectionInfo struct {
	Name        string `json:"name,omitempty"`
	PostmanID   string `json:"_postman_id,omitempty"`
	Description string `json:"description,omitempty"`
	Schema      string `json:"schema,omitempty"`
}

type FolderItem struct {
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Item        []APIItem `json:"item,omitempty"`
}

type APIItem struct {
	Name    string  `json:"name,omitempty"`
	Event   []Event `json:"event,omitempty"`
	Request Request `json:"request,omitempty"`
}

type Event struct {
	Listen string `json:"listen"`
	Script Script `json:"script"`
}

type Script struct {
	Type string   `json:"type,omitempty"`
	Exec []string `json:"exec,omitmpety"`
}

type Request struct {
	URL         URL         `json:"url,omitempty"`
	Method      string      `json:"method,omitempty"`
	Header      []Header    `json:"header,omitempty"`
	Body        RequestBody `json:"body,omitempty"`
	Description string      `json:"description,omitempty"`
}

type URL struct {
	Raw      string            `json:"raw,omitempty"`
	Protocol string            `json:"protocol,omitempty"`
	Auth     map[string]string `json:"auth"`
	Host     []string          `json:"host,omitempty"`
	Path     []string          `json:"path,omitempty"`
	Variable []URLVariable     `json:"variable,omitempty"`
}

type URLVariable struct {
	Value interface{} `json:"value,omitempty"`
	ID    string      `json:"id,omitempty"`
}

func NewURL(rawURL string) URL {
	rawURL = strings.TrimSpace(rawURL)
	url := URL{Raw: rawURL, Variable: []URLVariable{}}
	rx := regexp.MustCompile(`^([a-z]+)://([^/]+)/(.*)$`)
	rs := rx.FindAllStringSubmatch(rawURL, -1)

	if len(rs) > 0 {
		for _, m := range rs {
			url.Protocol = m[1]
			hostname := m[2]
			path := m[3]
			hostnameParts := strings.Split(hostname, ".")
			url.Host = hostnameParts

			pathParts := strings.Split(path, "/")
			url.Path = pathParts
		}
	}

	return url
}

func (url *URL) AddVariable(key string, value interface{}) {
	variable := URLVariable{ID: key, Value: value}
	url.Variable = append(url.Variable, variable)
}

type Header struct {
	Key         string `json:"key,omitempty"`
	Value       string `json:"value,omitempty"`
	Description string `json:"description,omitempty"`
}

type RequestBody struct {
	Mode       string            `json:"mode,omitempty"`
	URLEncoded []URLEncodedParam `json:"urlencoded,omitempty"`
}

type URLEncodedParam struct {
	Key     string `json:"key,omitempty"`
	Value   string `json:"value,omitempty"`
	Type    string `json:"type,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
}
