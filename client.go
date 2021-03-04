package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ernesto-jimenez/httplogger"
	"io"
	"io/ioutil"
	"net/http"
	"net/textproto"
	"time"
)

type HttpHeader map[string][]string

const (
	MethodGet    = "GET"
	MethodPost   = "POST"
	MethodPut    = "PUT"
	MethodDelete = "DELETE"
)

func (h HttpHeader) Set(key, value string) {
	textproto.MIMEHeader(h).Set(key, value)
}

func (h HttpHeader) Get(key string) string {
	return textproto.MIMEHeader(h).Get(key)
}

type Client struct {
	RetryMax    int // Maximum number of retries
	Timeout     time.Duration
	ShowHttpLog bool
}

func (client Client) GetForJson(url string, header HttpHeader, responseObject interface{}) error {
	return client.doForJson(MethodGet, url, header, nil, responseObject)
}

func (client Client) GetForJsonWithRequestObject(url string, header HttpHeader, requestObject interface{}, responseObject interface{}) error {
	return client.doForJson(MethodGet, url, header, requestObject, responseObject)
}

func (client Client) PostForJson(url string, header HttpHeader, requestObject interface{}) error {
	return client.doForJson(MethodPost, url, header, requestObject, nil)
}

func (client Client) PostForJsonWithResponseObject(url string, header HttpHeader, requestObject interface{}, responseObject interface{}) error {
	return client.doForJson(MethodPost, url, header, requestObject, responseObject)
}

func (client Client) PutForJson(url string, header HttpHeader, requestObject interface{}) error {
	return client.doForJson(MethodPut, url, header, requestObject, nil)
}

func (client Client) DeleteForJson(url string, header HttpHeader, requestObject interface{}) error {
	return client.doForJson(MethodDelete, url, header, requestObject, nil)
}

func (client Client) doForJson(method, url string, header HttpHeader, requestObject interface{}, responseObject interface{}) error {
	if header == nil {
		header = HttpHeader{}
		header.Set("Content-Type", "application/json;charset=UTF-8")
	}

	if len(header.Get("Content-Type")) == 0 {
		header.Set("Content-Type", "application/json;charset=UTF-8")
	}

	var requestBody io.Reader
	if requestObject != nil {
		marshal, err := json.Marshal(requestObject)
		if err != nil {
			return err
		}
		requestBody = bytes.NewBuffer(marshal)
	}

	response, err := client.Do(method, url, header, requestBody)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if responseObject != nil {
		b, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(b, &responseObject); err != nil {
			return err
		}
	}

	return nil
}

func (Client) Do(method, url string, header HttpHeader, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header(header)
	client := &http.Client{
		Transport: httplogger.NewLoggedTransport(http.DefaultTransport, newLogger()),
	}
	response, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	// TODO 200번대로 확장해야...
	if response.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		return nil, errors.New(fmt.Sprintf("error : %v, %v", response.StatusCode, string(b)))
	}

	return response, nil
}
