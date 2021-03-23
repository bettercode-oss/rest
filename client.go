package rest

import (
	"bytes"
	"encoding/json"
	"github.com/avast/retry-go"
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

var (
	defaultAttempts = uint(1)
	defaultDelay    = 2 * time.Second
)

type Client struct {
	RetryMax    uint
	RetryDelay  time.Duration
	Timeout     time.Duration
	ShowHttpLog bool
}

func (c Client) GetForJson(url string, header HttpHeader, responseObject interface{}) error {
	return c.doForJson(MethodGet, url, header, nil, responseObject)
}

func (c Client) GetForJsonWithRequestObject(url string, header HttpHeader, requestObject interface{}, responseObject interface{}) error {
	return c.doForJson(MethodGet, url, header, requestObject, responseObject)
}

func (c Client) PostForJson(url string, header HttpHeader, requestObject interface{}) error {
	return c.doForJson(MethodPost, url, header, requestObject, nil)
}

func (c Client) PostForJsonWithResponseObject(url string, header HttpHeader, requestObject interface{}, responseObject interface{}) error {
	return c.doForJson(MethodPost, url, header, requestObject, responseObject)
}

func (c Client) PutForJson(url string, header HttpHeader, requestObject interface{}) error {
	return c.doForJson(MethodPut, url, header, requestObject, nil)
}

func (c Client) DeleteForJson(url string, header HttpHeader, requestObject interface{}) error {
	return c.doForJson(MethodDelete, url, header, requestObject, nil)
}

func (c Client) doForJson(method, url string, header HttpHeader, requestObject interface{}, responseObject interface{}) error {
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

	response, err := c.Do(method, url, header, requestBody)
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

func (c Client) Do(method, url string, header HttpHeader, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header(header)
	client := &http.Client{
		Timeout:   c.Timeout,
		Transport: httplogger.NewLoggedTransport(http.DefaultTransport, newLogger()),
	}

	attempts := defaultAttempts
	delay := defaultDelay

	if c.RetryMax > 0 {
		attempts = c.RetryMax
	}

	if c.RetryDelay > 0 {
		delay = c.RetryDelay
	}

	var res *http.Response
	if err := retry.Do(
		func() error {
			response, err := client.Do(req)
			if err != nil {
				return err
			}

			if response.StatusCode >= http.StatusOK && response.StatusCode <= http.StatusIMUsed {
				res = response
				return nil
			}

			b, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return err
			}

			return &HttpServerError{
				StatusCode: response.StatusCode,
				Body:       string(b),
			}
		},
		retry.RetryIf(func(err error) bool {
			switch e := err.(type) {
			case *HttpServerError:
				if e.StatusCode >= http.StatusInternalServerError && e.StatusCode <= http.StatusNetworkAuthenticationRequired {
					return true
				}
				return false
			default:
				return false
			}

			return true
		}),
		retry.Attempts(attempts),
		retry.Delay(delay)); err != nil {
		return nil, err
	}

	return res, nil
}
