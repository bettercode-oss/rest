# Golang Rest Client

Golang에서 기본적으로 제공하는 [http 패키지](https://golang.org/pkg/net/http/) 는 JSON을 기반으로 하는 REST API를 호출할 때 
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
* REST 호출 시 자동 JSON 변환(마샬링/언마샬링) 제공
* HTTP Request & Response 로그 지원
* HTTP Request Timeout 지원
* HTTP Request Retry 지원

## 사용법
### 설치
```shell
go get github.com/bettercode-oss/rest
```

### REST API 호출 예시

Rest Client는 기본적으로 HTTP Request 헤더에 `Content-Type`을 `application/json;charset=UTF-8`로 설정한다.

#### GET

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

err := client.
	Request().
	SetResult(&responseObject).
	Get("http://example.com")
```

#### POST

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

err := client.
	Request().
	SetBody(requestObject).
	Post("http://example.com")
```

#### PUT

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

err := client.
	Request().
	SetBody(requestObject).
	Put("http://example.com")
```

#### DELETE

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

err := client.
	Request().
	SetBody(requestObject).
	Delete("http://example.com")
```

#### HTTP 헤더 추가 하기

Rest Client는 기본적으로 헤더에 `Content-Type`에 `application/json;charset=UTF-8`을 추가한다.
이외에 헤더를 추가하고 싶다면 아래 처럼 추가한다.
```go
client := rest.Client{}
responseObject := struct {
  Id   string `json:"id"`
  Name string `json:"name"`
  Age  int    `json:"age"`
}{}

err := client.
	Request().
        SetHeader("Accept", "application/json").
	SetHeader("Authorization", "a9ace025c90c0da2161075da6ddd3492a2fca776").
	SetResult(&responseObject).
	Get("http://example.com")
```

### HTTP Request/Response 로깅(Logging)
Rest Client 생성할 때 `ShowHttpLog`를 `true`로 설정한다.
```go
client := rest.Client{ShowHttpLog: true}
```
아래와 같이 로그를 확인할 수 있다.
```
bettercode-oss/rest - 2021/03/11 07:06:36 Request method=GET url=http://localhost:51716 header=map[Authorization:[a9ace025c90c0da2161075da6ddd3492a2fca776] Content-Type:[application/json;charset=UTF-8]] body=
bettercode-oss/rest - 2021/03/11 07:06:36 Response method=GET url=http://localhost:51716 status=200 durationMs=4 header=map[Content-Length:[89] Content-Type:[text/plain; charset=utf-8] Date:[Wed, 10 Mar 2021 22:06:36 GMT]] 
body={
    "id": "gigamadness@gmail.com",
    "name": "Yoo Young-mo",
    "age": 20
}
```

### HTTP Request Timeout 설정
Rest Client 생성할 때 `Timeout`을 지정한다.
```go
client := rest.Client{
  Timeout:     10 * time.Second
}
```

### HTTP Request Retry
Rest Client 생성할 때 재시도 최대 횟수(`RetryMax`)를 지정한다.
추가로 재시도 지연시간(`RetryDelay`)을 지정할 수 있다.
```go
client := rest.Client{
  RetryMax:    5,
  RetryDelay:  1 * time.Second,
}
```
위와 같이 설정하면 HTTP 응답 코드가 500번대(Server Error)인 경우 다시 HTTP 요청한다.
최대 5번 하게 되며 시도 사이의 지연 시간은 1초이다.

### Error 처리
HTTP Status 에 따라 에러 처리가 필요한 경우 아래와 같이 확인할 수 있다.
주의. HttpServerError 만 반환하는 것은 아니다. HTTP Request 과정에서 발생한 error 만  HttpServerError를 반환한다.
```go
err := client.
		Request().
		SetResult(&responseObject).
		Get("http://example.com/err-url")

if httpErr, ok := err.(*HttpServerError); ok {
	// ...  		
}
```

### SSL 인증서 오류(x509) 무시 옵션
`InsecureSkipVerify` 를 true로 설정한다.
```go
client := rest.Client{
  // ...
  InsecureSkipVerify: true,
}
```
