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

func TestClient_GetForJson(t *testing.T) {
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
		RetryMax:    10,
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
	err := client.GetForJson(fmt.Sprintf("http://localhost:%v", serverPort), header, &responseObject)

	// then
	assert.Nil(t, err)
	assert.Equal(t, "gigamadness@gmail.com", responseObject.Id)
	assert.Equal(t, "Yoo Young-mo", responseObject.Name)
	assert.Equal(t, 20, responseObject.Age)
}

func TestClient_GetForJsonWithRequestObject(t *testing.T) {
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
		RetryMax:    10,
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
	err := client.GetForJsonWithRequestObject(fmt.Sprintf("http://localhost:%v", serverPort), nil, requestObject, &responseObject)

	// then
	assert.Nil(t, err)
	assert.Equal(t, "gigamadness@gmail.com", responseObject.Id)
	assert.Equal(t, "Yoo Young-mo", responseObject.Name)
	assert.Equal(t, 20, responseObject.Age)
}

func TestClient_PostForJson(t *testing.T) {
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
		RetryMax:    10,
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
	err := client.PostForJson(fmt.Sprintf("http://localhost:%v", serverPort), nil, requestObject)

	// then
	assert.Nil(t, err)
}

func TestClient_PostForJsonWithResponseObject(t *testing.T) {
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
		RetryMax:    10,
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
	err := client.PostForJsonWithResponseObject(fmt.Sprintf("http://localhost:%v", serverPort), nil, requestObject, &responseObject)

	// then
	assert.Nil(t, err)
	assert.Equal(t, "gigamadness@gmail.com", responseObject.Id)
	assert.Equal(t, "Yoo Young-mo", responseObject.Name)
	assert.Equal(t, 20, responseObject.Age)
}

func TestClient_PutForJson(t *testing.T) {
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
		RetryMax:    10,
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
	err := client.PutForJson(fmt.Sprintf("http://localhost:%v", serverPort), nil, requestObject)

	// then
	assert.Nil(t, err)
}

func TestClient_DeleteForJson(t *testing.T) {
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
		RetryMax:    10,
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
	err := client.DeleteForJson(fmt.Sprintf("http://localhost:%v", serverPort), nil, requestObject)

	// then
	assert.Nil(t, err)
}

func TestClient_GetForJson_Error_Timeout(t *testing.T) {
	// setUp WebServer Fixture
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			time.Sleep(20 * time.Second)

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
	err := client.GetForJson(fmt.Sprintf("http://localhost:%v", serverPort), header, &responseObject)

	// then
	assert.NotNil(t, err)
}
