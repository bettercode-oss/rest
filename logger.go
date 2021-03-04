package rest

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type httpLogger struct {
	log *log.Logger
}

func newLogger() *httpLogger {
	return &httpLogger{
		log: log.New(os.Stderr, "bettercode-oss/http - ", log.LstdFlags),
	}
}

func (logger *httpLogger) LogRequest(req *http.Request) {
	var body string

	if req.Body != nil {
		// Read the content
		rawBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			logger.log.Println(err)
		}

		// Restore the io.ReadCloser to it's original state
		req.Body = ioutil.NopCloser(bytes.NewBuffer(rawBody))
		body = string(rawBody)
	}

	logger.log.Printf(
		"Request method=%s url=%s header=%s body=%s",
		req.Method,
		req.URL.String(),
		req.Header,
		body,
	)
}

func (logger *httpLogger) LogResponse(req *http.Request, res *http.Response, err error, duration time.Duration) {
	duration /= time.Millisecond
	if err != nil {
		logger.log.Println(err)
	} else {
		var body string
		if res.Body != nil {
			// Read the content
			rawBody, err := ioutil.ReadAll(res.Body)
			if err != nil {
				logger.log.Println(err)
			}

			// Restore the io.ReadCloser to it's original state
			res.Body = ioutil.NopCloser(bytes.NewBuffer(rawBody))
			body = string(rawBody)
		}

		logger.log.Printf(
			"Response method=%s url=%s status=%d durationMs=%d header=%s body=%s",
			req.Method,
			req.URL.String(),
			res.StatusCode,
			duration,
			res.Header,
			body,
		)
	}
}
