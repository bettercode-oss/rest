package rest

import "fmt"

type HttpServerError struct {
	StatusCode int
	Body       string
}

func (e *HttpServerError) Error() string {
	return fmt.Sprintf("bettercode-oss/rest: http server error - status code: %v; body: %v", e.StatusCode, e.Body)
}
