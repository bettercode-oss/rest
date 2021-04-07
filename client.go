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

func (c *Client) doForJson(r *Request) error {
	if r.header == nil {
		r.header = HttpHeader{}
		r.header.Set("Content-Type", "application/json;charset=UTF-8")
	}

	if len(r.header.Get("Content-Type")) == 0 {
		r.header.Set("Content-Type", "application/json;charset=UTF-8")
	}

	var requestBody io.Reader
	if r.body != nil {
		marshal, err := json.Marshal(r.body)
		if err != nil {
			return err
		}
		requestBody = bytes.NewBuffer(marshal)
	}

	response, err := c.do(r.method, r.url, r.header, requestBody)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if r.result != nil {
		b, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(b, &r.result); err != nil {
			return err
		}
	}

	return nil
}

func (c Client) do(method, url string, header HttpHeader, body io.Reader) (*http.Response, error) {
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

func (c *Client) Request() *Request {
	return &Request{client: c}
}

type Request struct {
	client *Client
	url    string
	method string
	header HttpHeader
	body   interface{}
	result interface{}
}

func (r *Request) SetHeader(key, value string) *Request {
	if r.header == nil {
		r.header = HttpHeader{}
	}

	r.header.Set(key, value)
	return r
}

func (r *Request) SetBody(body interface{}) *Request {
	r.body = body
	return r
}

func (r *Request) SetResult(result interface{}) *Request {
	r.result = result
	return r
}

func (r *Request) Get(url string) error {
	r.method = MethodGet
	r.url = url
	return r.client.doForJson(r)
}

func (r *Request) Post(url string) error {
	r.method = MethodPost
	r.url = url
	return r.client.doForJson(r)
}

func (r *Request) Delete(url string) error {
	r.method = MethodDelete
	r.url = url
	return r.client.doForJson(r)
}

func (r *Request) Put(url string) error {
	r.method = MethodPut
	r.url = url
	return r.client.doForJson(r)
}
