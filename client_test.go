package rest

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_Get(t *testing.T) {
	// setUp WebServer Fixture
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json")
			responseBody := `{
				"id": "gigamadness@gmail.com",
				"name": "Yoo Young-mo",
        		"age": 20
			}`
			w.Write([]byte(responseBody))
		} else {
			w.WriteHeader(404)
		}
	}))
	defer server.Close()
	serverPort := server.Listener.Addr().(*net.TCPAddr).Port

	// given
	client := Client{
		Timeout:     10 * time.Second,
		ShowHttpLog: true,
	}

	responseObject := struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{}

	// when
	err := client.
		Request().
		SetHeader("Authorization", "a9ace025c90c0da2161075da6ddd3492a2fca776").
		SetResult(&responseObject).
		Get(fmt.Sprintf("http://localhost:%v", serverPort))

	// then
	assert.Nil(t, err)
	assert.Equal(t, "gigamadness@gmail.com", responseObject.Id)
	assert.Equal(t, "Yoo Young-mo", responseObject.Name)
	assert.Equal(t, 20, responseObject.Age)
}

func TestClient_Get_WithRequestBody(t *testing.T) {
	// setUp WebServer Fixture
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			requestBody := map[string]interface{}{}
			b, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(b, &requestBody)

			if len(requestBody["id"].(string)) > 0 {
				w.WriteHeader(200)
				w.Header().Set("Content-Type", "application/json")
				responseBody := fmt.Sprintf(`{
					"id": "%s",
					"name": "Yoo Young-mo",
					"age": 20
				}`, requestBody["id"].(string))
				w.Write([]byte(responseBody))
			} else {
				w.WriteHeader(400)
			}
		} else {
			w.WriteHeader(404)
		}
	}))
	defer server.Close()
	serverPort := server.Listener.Addr().(*net.TCPAddr).Port

	// given
	client := Client{
		Timeout:     10 * time.Second,
		ShowHttpLog: true,
	}

	requestObject := struct {
		Id string `json:"id"`
	}{
		Id: "gigamadness@gmail.com",
	}

	responseObject := struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{}

	// when
	err := client.
		Request().
		SetBody(requestObject).
		SetResult(&responseObject).
		Get(fmt.Sprintf("http://localhost:%v", serverPort))

	// then
	assert.Nil(t, err)
	assert.Equal(t, "gigamadness@gmail.com", responseObject.Id)
	assert.Equal(t, "Yoo Young-mo", responseObject.Name)
	assert.Equal(t, 20, responseObject.Age)
}

func TestClient_Post(t *testing.T) {
	// setUp WebServer Fixture
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			w.WriteHeader(200)
		} else {
			w.WriteHeader(404)
		}
	}))
	defer server.Close()
	serverPort := server.Listener.Addr().(*net.TCPAddr).Port

	// given
	client := Client{
		Timeout:     10 * time.Second,
		ShowHttpLog: true,
	}

	requestObject := struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Id:   "gigamadness@gmail.com",
		Name: "Yoo Young-mo",
		Age:  20,
	}

	// when
	err := client.
		Request().
		SetBody(requestObject).
		Post(fmt.Sprintf("http://localhost:%v", serverPort))

	// then
	assert.Nil(t, err)
}

func TestClient_Post_WithResponseObject(t *testing.T) {
	// setUp WebServer Fixture
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			w.WriteHeader(200)
			requestBody := map[string]interface{}{}
			b, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(b, &requestBody)

			if len(requestBody["id"].(string)) > 0 {
				w.WriteHeader(200)
				w.Header().Set("Content-Type", "application/json")
				responseBody := fmt.Sprintf(`{
					"id": "%s",
					"name": "Yoo Young-mo",
					"age": 20
				}`, requestBody["id"].(string))
				w.Write([]byte(responseBody))
			} else {
				w.WriteHeader(400)
			}
		} else {
			w.WriteHeader(404)
		}
	}))
	defer server.Close()
	serverPort := server.Listener.Addr().(*net.TCPAddr).Port

	// given
	client := Client{
		Timeout:     10 * time.Second,
		ShowHttpLog: true,
	}

	requestObject := struct {
		Id string `json:"id"`
	}{
		Id: "gigamadness@gmail.com",
	}

	responseObject := struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{}

	// when
	err := client.
		Request().
		SetBody(requestObject).
		SetResult(&responseObject).
		Post(fmt.Sprintf("http://localhost:%v", serverPort))

	// then
	assert.Nil(t, err)
	assert.Equal(t, "gigamadness@gmail.com", responseObject.Id)
	assert.Equal(t, "Yoo Young-mo", responseObject.Name)
	assert.Equal(t, 20, responseObject.Age)
}

func TestClient_Put(t *testing.T) {
	// setUp WebServer Fixture
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			w.WriteHeader(200)
		} else {
			w.WriteHeader(404)
		}
	}))
	defer server.Close()
	serverPort := server.Listener.Addr().(*net.TCPAddr).Port

	// given
	client := Client{
		Timeout:     10 * time.Second,
		ShowHttpLog: true,
	}

	requestObject := struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Id:   "gigamadness@gmail.com",
		Name: "Yoo Young-mo",
		Age:  20,
	}

	// when
	err := client.
		Request().
		SetBody(requestObject).
		Put(fmt.Sprintf("http://localhost:%v", serverPort))

	// then
	assert.Nil(t, err)
}

func TestClient_Delete(t *testing.T) {
	// setUp WebServer Fixture
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			w.WriteHeader(200)
		} else {
			w.WriteHeader(404)
		}
	}))
	defer server.Close()
	serverPort := server.Listener.Addr().(*net.TCPAddr).Port

	// given
	client := Client{
		Timeout:     10 * time.Second,
		ShowHttpLog: true,
	}

	requestObject := struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Id:   "gigamadness@gmail.com",
		Name: "Yoo Young-mo",
		Age:  20,
	}

	// when
	err := client.Request().SetBody(requestObject).Delete(fmt.Sprintf("http://localhost:%v", serverPort))

	// then
	assert.Nil(t, err)
}

func TestClient_Error_Timeout(t *testing.T) {
	// setUp WebServer Fixture
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			time.Sleep(10 * time.Second)

			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json")
			responseBody := `{
				"id": "gigamadness@gmail.com",
				"name": "Yoo Young-mo",
       		"age": 20
			}`
			w.Write([]byte(responseBody))
		} else {
			w.WriteHeader(404)
		}
	}))
	defer server.Close()
	serverPort := server.Listener.Addr().(*net.TCPAddr).Port

	// given
	client := Client{
		Timeout:     5 * time.Second,
		ShowHttpLog: true,
	}

	header := HttpHeader{}
	header.Set("Authorization", "a9ace025c90c0da2161075da6ddd3492a2fca776")

	responseObject := struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{}

	// when
	err := client.Request().
		SetHeader("Authorization", "a9ace025c90c0da2161075da6ddd3492a2fca776").
		SetResult(&responseObject).
		Get(fmt.Sprintf("http://localhost:%v", serverPort))

	// then
	assert.NotNil(t, err)
}

func TestClient_Retry_Response_Http_code_500_Eventually_Success(t *testing.T) {
	// setUp WebServer Fixture
	serverErrorCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if serverErrorCount < 4 {
				w.WriteHeader(500)
				serverErrorCount++
			} else {
				w.WriteHeader(200)
				w.Header().Set("Content-Type", "application/json")
				responseBody := `{
					"id": "gigamadness@gmail.com",
					"name": "Yoo Young-mo",
					"age": 20
				}`
				w.Write([]byte(responseBody))
			}
		} else {
			w.WriteHeader(404)
		}
	}))
	defer server.Close()
	serverPort := server.Listener.Addr().(*net.TCPAddr).Port

	// given
	client := Client{
		Timeout:     10 * time.Second,
		ShowHttpLog: true,
		RetryMax:    5,
		RetryDelay:  1 * time.Second,
	}

	responseObject := struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{}

	// when
	err := client.
		Request().
		SetResult(&responseObject).
		Get(fmt.Sprintf("http://localhost:%v", serverPort))

	// then
	assert.Nil(t, err)
	assert.Equal(t, "gigamadness@gmail.com", responseObject.Id)
	assert.Equal(t, "Yoo Young-mo", responseObject.Name)
	assert.Equal(t, 20, responseObject.Age)
}

func TestClient_Retry_Response_Http_code_500_Eventually_Failed(t *testing.T) {
	// setUp WebServer Fixture
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.WriteHeader(500)
		} else {
			w.WriteHeader(404)
		}
	}))
	defer server.Close()
	serverPort := server.Listener.Addr().(*net.TCPAddr).Port

	// given
	client := Client{
		Timeout:     10 * time.Second,
		ShowHttpLog: true,
		RetryMax:    5,
		RetryDelay:  1 * time.Second,
	}

	responseObject := struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{}

	// when
	err := client.Request().SetResult(&responseObject).Get(fmt.Sprintf("http://localhost:%v", serverPort))

	// then
	assert.NotNil(t, err)
	if httpErr, ok := err.(*HttpServerError); ok {
		assert.NotNil(t, httpErr)
		assert.Equal(t, http.StatusInternalServerError, httpErr.StatusCode)
	} else {
		t.Fail()
	}
}

func TestClient_Retry_Response_Http_code_400(t *testing.T) {
	// setUp WebServer Fixture
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(404)
		}
	}))
	defer server.Close()
	serverPort := server.Listener.Addr().(*net.TCPAddr).Port

	// given
	client := Client{
		Timeout:     10 * time.Second,
		ShowHttpLog: true,
		RetryMax:    5,
		RetryDelay:  1 * time.Second,
	}

	responseObject := struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{}

	// when
	err := client.
		Request().
		SetResult(&responseObject).
		Get(fmt.Sprintf("http://localhost:%v", serverPort))

	// then
	assert.NotNil(t, err)
	if httpErr, ok := err.(*HttpServerError); ok {
		assert.NotNil(t, httpErr)
		assert.Equal(t, http.StatusBadRequest, httpErr.StatusCode)
	} else {
		t.Fail()
	}
}

func TestClient_SSL_x509_error(t *testing.T) {
	// x509 error test server
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.WriteHeader(200)
		} else {
			w.WriteHeader(404)
		}
	}))

	defer server.Close()
	serverPort := server.Listener.Addr().(*net.TCPAddr).Port

	client := Client{
		Timeout:            10 * time.Second,
		ShowHttpLog:        true,
		RetryMax:           5,
		RetryDelay:         1 * time.Second,
		InsecureSkipVerify: true,
	}

	// when
	err := client.
		Request().
		Get(fmt.Sprintf("https://localhost:%v", serverPort))

	assert.Nil(t, err)
}
