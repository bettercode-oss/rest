# Golang Rest Client

## 배경
Golang에서 기본적으로 제공하는 [http 패키지](https://golang.org/pkg/net/http/) 로 JSON을 기반으로 하는 REST API를 호출할 때 
HTTP Request/Response Body를 JSON으로 변환(json.Marshal/Unmarshal)해야 하는 불편함이 있다.

아래는 GET/POST 방식 호출 예시이다.
```golang
// HTTP POST 예시
requestBody := map[string]interface{}
requestBody["name"] = "Yoo Young-mo"
requestBody["email"] = "gigamadness@gmail.com"

pbytes, _ := json.Marshal(requestBody)
buff := bytes.NewBuffer(pbytes)
resp, err := http.Post("http://example.com", "application/json", buff)
// ...
```
GET 방식으로 호출할 때 `http.Get` 함수는 HTTP Header를 지원하지 않기 때문에 아래와 같이 http.Client를 따로 만들어서 사용해야 한다.
```golang
// HTTP GET 예시
req, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
req.Header.Set("Content-Type", "application/json;charset=UTF-8")

client := &http.Client{}
resp, err := client.Do(req)
defer resp.Body.Close()

bytes, _ := ioutil.ReadAll(resp.Body)
// ...
var responseBody = map[string]interface{}{}
json.Unmarshal(bytes, &responseBody)
```

## 제공하는 기능
* [X] REST API 호출 Client 제공
* [X] HTTP Request & Response 로그 지원
* [ ] HTTP Request Timeout 지원(HTTP Request 시 특정 시간이 지나도 응답이 지연되는 경우 강제로 커넥션을 끊게 만든다)
* [ ] HTTP Request Retry 지원(실패시 최대 재시도 횟수를 지정할 수 있게 만든다.)

## 사용법
### 설치
```shell
go get github.com/bettercode-oss/rest
```

### REST API 호출 예시

* GET

```go
import (
  "github.com/bettercode-oss/rest"
)
// ...
client := rest.Client{}
responseObject := struct {
  Id   string `json:"id"`
  Name string `json:"name"`
  Age  int    `json:"age"`
}{}

err := client.GetForJson("http://example.com", nil, &responseObject)
```

* POST

```go
import (
  "github.com/bettercode-oss/rest"
)
// ...
client := rest.Client{}
requestObject := struct {
  Id   string `json:"id"`
  Name string `json:"name"`
  Age  int    `json:"age"`
}{
  Id:   "gigamadness@gmail.com",
  Name: "Yoo Young-mo",
  Age:  20,
}

err := client.PostForJson("http://example.com", nil, requestObject)
```

* PUT

```go
import (
  "github.com/bettercode-oss/rest"
)
// ...
client := rest.Client{}
requestObject := struct {
  Id   string `json:"id"`
  Name string `json:"name"`
  Age  int    `json:"age"`
}{
  Id:   "gigamadness@gmail.com",
  Name: "Yoo Young-mo",
  Age:  20,
}

err := client.PutForJson("http://example.com", nil, requestObject)
```

* DELETE

```go
import (
  "github.com/bettercode-oss/rest"
)
// ...
client := rest.Client{}
requestObject := struct {
  Id   string `json:"id"`
  Name string `json:"name"`
  Age  int    `json:"age"`
}{
  Id:   "gigamadness@gmail.com",
  Name: "Yoo Young-mo",
  Age:  20,
}

err := client.DeleteForJson("http://example.com", nil, requestObject)
```

* HTTP 헤더 추가 하기

Rest Client는 기본적으로 헤더에 `Content-Type`에 `application/json;charset=UTF-8`을 추가한다.
이외에 헤더를 추가하고 싶다면 아래 처럼 추가한다.
```go
client := rest.Client{}
header := rest.HttpHeader{}
header.Set("Authorization", "a9ace025c90c0da2161075da6ddd3492a2fca776")
responseObject := struct {
  Id   string `json:"id"`
  Name string `json:"name"`
  Age  int    `json:"age"`
}{}
err := client.GetForJson("http://example.com", header, &responseObject)
```

### HTTP Request/Response 로깅(Logging)
Rest Client 생성 시 `ShowHttpLog`를 `true`로 설정한다.
```go
client := rest.Client{ShowHttpLog: true}
```
아래와 같이 로그를 확인할 수 있다.
```
bettercode-oss/rest - 2021/03/11 07:06:36 Request method=GET url=http://localhost:51716 header=map[Authorization:[a9ace025c90c0da2161075da6ddd3492a2fca776] Content-Type:[application/json;charset=UTF-8]] body=
bettercode-oss/rest - 2021/03/11 07:06:36 Response method=GET url=http://localhost:51716 status=200 durationMs=4 header=map[Content-Length:[89] Content-Type:[text/plain; charset=utf-8] Date:[Wed, 10 Mar 2021 22:06:36 GMT]] body={
				"id": "gigamadness@gmail.com",
				"name": "Yoo Young-mo",
        		"age": 20
			}
```
### HTTP Request Timeout 설정
TO-DO

### HTTP Request Retry
TO-DO
