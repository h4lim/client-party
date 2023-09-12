package party

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	MIMEJSON              = "application/json"
	MIMEHTML              = "text/html"
	MIMEXML               = "application/xml"
	MIMEXML2              = "text/xml"
	MIMEPlain             = "text/plain"
	MIMEPOSTForm          = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm = "multipart/form-data"
	MIMEPROTOBUF          = "application/x-protobuf"
	MIMEMSGPACK           = "application/x-msgpack"
	MIMEMSGPACK2          = "application/msgpack"
	MIMEYAML              = "application/x-yaml"
)

const (
	MethodGet     = "GET"
	MethodHead    = "HEAD"
	MethodPost    = "POST"
	MethodPut     = "PUT"
	MethodPatch   = "PATCH" // RFC 5789
	MethodDelete  = "DELETE"
	MethodConnect = "CONNECT"
	MethodOptions = "OPTIONS"
	MethodTrace   = "TRACE"
)

type ClientParty struct {
	HttpMethod  string
	URL         string
	Header      map[string]string
	QueryParam  *map[string]string
	BaseAuth    *map[string]string
	RequestBody *[]byte
	Writer      *multipart.Writer
	HttpClient  http.Client
}

type Response struct {
	HttpCode     int
	ResponseBody string
}

type ClientPartyBuilder struct {
	ClientParty ClientParty
}

type IClientPartyBuilder interface {
	SetHeader(contentType string, header map[string]string) ClientPartyBuilder
	SetQueryParam(map[string]string) ClientPartyBuilder
	SetBaseAuth(username string, password string) ClientPartyBuilder
	SetRequestBody(requestBody interface{}) (*ClientPartyBuilder, *error)
	SetRequestBodyStr(requestBody string) ClientPartyBuilder
	SetFormData(mapFile map[string]string, mapText map[string]string) (*ClientPartyBuilder, *error)
	HitClient() (*Response, *error)
}

func NewClientParty(httpMethod string, url string) IClientPartyBuilder {

	return ClientPartyBuilder{ClientParty: ClientParty{
		HttpMethod: httpMethod,
		URL:        url,
	}}
}

func (c ClientPartyBuilder) SetHeader(contentType string, header map[string]string) ClientPartyBuilder {

	if contentType != "" {
		header["Content-Type"] = contentType
	}

	c.ClientParty.Header = header

	return c
}

func (c ClientPartyBuilder) SetQueryParam(queryParam map[string]string) ClientPartyBuilder {

	c.ClientParty.QueryParam = &queryParam

	return c
}

func (c ClientPartyBuilder) SetBaseAuth(username string, password string) ClientPartyBuilder {

	mapBaseAuth := make(map[string]string)
	mapBaseAuth[username] = password

	c.ClientParty.BaseAuth = &mapBaseAuth

	return c
}

func (c ClientPartyBuilder) SetRequestBody(requestBody interface{}) (*ClientPartyBuilder, *error) {

	contentType := c.ClientParty.Header["Content-Type"]
	if contentType == "" {

		byteRequest, err := json.Marshal(requestBody)
		if err != nil {
			return nil, &err
		}
		c.ClientParty.RequestBody = &byteRequest
	}

	if contentType == MIMEJSON {

		byteRequest, err := json.Marshal(requestBody)
		if err != nil {
			return nil, &err
		}
		c.ClientParty.RequestBody = &byteRequest
	}

	if contentType == MIMEXML || contentType == MIMEXML2 {

		byteRequest, err := xml.Marshal(requestBody)
		if err != nil {
			return nil, &err
		}

		c.ClientParty.RequestBody = &byteRequest
	}

	if contentType == MIMEPOSTForm {

		byteRequest, err := json.Marshal(requestBody)
		if err != nil {
			return nil, &err
		}

		mapJson := map[string]string{}
		if err := json.Unmarshal(byteRequest, &mapJson); err != nil {
			return nil, &err
		}

		data := url.Values{}
		for key, value := range mapJson {
			data.Set(key, value)
		}

		bytePostForm := []byte(data.Encode())
		c.ClientParty.RequestBody = &bytePostForm

	}

	if contentType == MIMEMultipartPOSTForm {

	}

	return &c, nil
}

func (c ClientPartyBuilder) SetRequestBodyStr(requestBody string) ClientPartyBuilder {
	byteRequestBody := []byte(requestBody)
	c.ClientParty.RequestBody = &byteRequestBody
	return c
}

func (c ClientPartyBuilder) SetFormData(mapFile map[string]string, mapText map[string]string) (*ClientPartyBuilder, *error) {

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for k, v := range mapFile {

		fw, err := writer.CreateFormFile(k, v)
		if err != nil {
			return nil, &err
		}

		file, err := os.Open(v)
		if err != nil {
			return nil, &err
		}

		if _, err := io.Copy(fw, file); err != nil {
			return nil, &err
		}

	}

	for k, v := range mapText {

		fw, err := writer.CreateFormField(k)
		if err != nil {
			return nil, &err
		}

		if _, err = io.Copy(fw, strings.NewReader(v)); err != nil {
			return nil, &err
		}

	}

	c.ClientParty.Writer = writer

	return &c, nil
}

func (c ClientPartyBuilder) HitClient() (*Response, *error) {

	var ioRequest io.Reader = nil
	if c.ClientParty.RequestBody != nil {
		ioRequest = bytes.NewReader(*c.ClientParty.RequestBody)
	}

	request, err := http.NewRequest(c.ClientParty.HttpMethod, c.ClientParty.URL, ioRequest)
	if err != nil {
		return nil, &err
	}

	if c.ClientParty.QueryParam != nil {
		q := request.URL.Query()
		for k, v := range *c.ClientParty.QueryParam {
			q.Add(k, v)
		}
		request.URL.RawQuery = q.Encode()
	}

	for k, v := range c.ClientParty.Header {
		if c.ClientParty.Writer != nil {
			request.Header.Set("Content-Type", c.ClientParty.Writer.FormDataContentType())
		} else {
			request.Header.Set(k, v)
		}
	}

	if c.ClientParty.BaseAuth != nil {
		for k, v := range *c.ClientParty.BaseAuth {
			request.SetBasicAuth(k, v)
		}
	}

	response, err := c.ClientParty.HttpClient.Do(request)
	if err != nil {
		return nil, &err
	}

	byteResult, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, &err
	}

	clientResponse := Response{
		HttpCode:     response.StatusCode,
		ResponseBody: string(byteResult),
	}

	return &clientResponse, nil

}
