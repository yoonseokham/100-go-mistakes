# 변수 섀도잉 (#1)

## 문제

`:=`를 if/else 블록 안에서 쓰면 외부 변수와 이름이 같은 **새 변수**가 생성된다. 외부 변수는 초기화되지 않은 채로 남는다.

```go
var client *http.Client
if tracing {
    client, err := createTracingClient()  // 새로운 client 변수 생성 (외부 변수 아님)
    // ...
}
// client는 여전히 nil!
```

## 해결책

### 해결책 1: `:=` 대신 `=` 사용

새 변수를 만드는 대신 기존 변수에 할당:

```go
var client *http.Client
var err error

if tracing {
    client, err = createTracingClient()  // 외부 변수에 할당
    // ...
}
```

### 해결책 2: 임시 변수 사용

블록 안에서 임시 변수를 쓰고 나중에 외부 변수에 할당:

```go
var client *http.Client

if tracing {
    c, err := createTracingClient()  // 임시 변수 c
    if err != nil {
        return nil, err
    }
    client = c  // 외부 변수에 할당
}
```

## 핵심 원칙

- 기존 변수에 값을 넣을 땐 `=` 사용
- 새 변수를 선언할 때만 `:=` 사용
- 섀도잉은 컴파일 에러가 아니라서 찾기 어려운 버그가 됨
