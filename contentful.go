package contentful

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

const baseURL = "https://cdn.contentful.com/spaces/"

type Client struct {
	AccessToken string
	Space       string
}

func New(token string, space string) *Client {
	return &Client{token, space}
}

type SystemProperties struct {
	Space       interface{}
	Type        string
	Id          string
	ContentType interface{}
	Revision    int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Locale struct {
	Code    string
	Default bool
	Name    string
}

func (client *Client) query(endpoint string, params map[string]string) (resp []byte, err error) {
	c := &http.Client{}
	req, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		return
	}

	req.Header.Set("Authorization", "Bearer "+client.AccessToken)
	res, err := c.Do(req)

	if err != nil {
		return
	}
	if res.StatusCode != http.StatusOK {
		return resp, errors.New(res.Status)
	}
	resp, err = ioutil.ReadAll(res.Body)
	return
}

// Spaces are containers for Content Types, Entries and Assets. API consumers,
// like mobile apps or websites, typically fetch data by getting Entries and
// Assets from one or more Spaces.
type Space struct {
	Sys     SystemProperties
	Name    string
	Locales []Locale
}

func (client *Client) getSpaceURL() string {
	return baseURL + client.Space
}

func (client *Client) GetSpace() (resp *Space, err error) {
	body, err := client.query(client.getSpaceURL(), nil)

	if err != nil {
		return
	}

	err = json.Unmarshal(body, &resp)

	return
}

// Content Types are schemas describing the shape of Entries. They mainly
// consist of a list of fields acting as a blueprint for Entries.
type ContentTypes struct {
	Sys   SystemProperties
	Total int64
	Skip  int64
	Limit int64
	Items []ContentType
}

type ContentType struct {
	Sys         SystemProperties
	Name        string
	Description string
	Fields      []Field
}

type Field struct {
	Id        string
	Name      string
	Type      string
	LinkType  string
	Items     FieldItemType
	Required  bool
	Localized bool
}

type FieldItemType struct {
	Type string
}

func (client *Client) getContentTypesURL() string {
	return client.getSpaceURL() + "/content_types"
}

func (client *Client) GetContentTypes() (resp *ContentTypes, err error) {
	body, err := client.query(client.getContentTypesURL(), nil)

	if err != nil {
		return
	}

	err = json.Unmarshal(body, &resp)

	return
}

func (client *Client) getContentTypeURL(id string) string {
	return client.getContentTypesURL() + "/" + id
}

func (client *Client) GetContentType(id string) (resp *ContentType, err error) {
	body, err := client.query(client.getContentTypeURL(id), nil)

	if err != nil {
		return
	}

	err = json.Unmarshal(body, &resp)

	return
}

// Entries represent textual content in a Space. An Entry's data adheres
// to a certain Content Type.
type Entries struct {
	Sys   SystemProperties
	Total int64
	Skip  int64
	Limit int64
	Items []Entry
}

type Entry struct {
	Sys    SystemProperties
	Fields interface{}
}

func (client *Client) getEntriesURL() string {
	return client.getSpaceURL() + "/entries"
}

func (client *Client) getEntries() (resp *Entries, err error) {
	body, err := client.query(client.getEntriesURL(), nil)

	if err != nil {
		return
	}

	err = json.Unmarshal(body, &resp)

	return
}

func (client *Client) getEntryURL(id string) string {
	return client.getEntriesURL() + "/" + id
}

func (client *Client) getEntry(id string) (resp *Entry, err error) {
	body, err := client.query(client.getEntryURL(id), nil)

	if err != nil {
		return
	}

	err = json.Unmarshal(body, &resp)

	return
}

type Assets struct {
	Sys   SystemProperties
	Total int64
	Skip  int64
	Limit int64
	Items []Asset
}

type Asset struct {
	Sys    SystemProperties
	Fields struct {
		Title string
		File  File
	}
}

type File struct {
	FileName    string
	ContentType string
	URL         string
	Details     interface{}
}

func (client *Client) getAssetsURL() string {
	return client.getSpaceURL() + "/assets"
}

func (client *Client) getAssets() (resp *Assets, err error) {
	body, err := client.query(client.getAssetsURL(), nil)

	if err != nil {
		return
	}

	err = json.Unmarshal(body, &resp)

	return
}

// func main() {
// 	client := new(Client)
// 	client.Space = "cfexampleapi"
// 	client.AccessToken = "b4c0n73n7fu1"

// 	// u := client.getEntriesURL()
// 	// fmt.Println("working", u)

// 	s, err := client.getAssets()
// 	fmt.Println("working", s, err)
// }
