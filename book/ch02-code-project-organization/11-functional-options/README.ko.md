# Functional Options 패턴을 사용하지 않는 실수 (#11)

## 한 줄 요약

선택적 설정값은 functional options 패턴으로 처리하면 API가 깔끔하고 확장하기 쉽다.

---

## 문제

생성자에 선택적 파라미터가 많을 때 흔히 쓰는 두 가지 방법은 모두 단점이 있다.

### Bad V1: 파라미터 폭발

```go
func NewServer(addr string, port int, timeout time.Duration, maxConn int) (*http.Server, error)
```

- 옵션이 추가될 때마다 시그니처가 바뀜 → 모든 호출부 수정 필요
- 관심 없는 값도 모두 전달해야 함

### Bad V2: Config 구조체

```go
func NewServer(addr string, cfg Config) (*http.Server, error)
```

- `Config{}`를 넘기는 게 어색하고 의미가 불명확
- "포트 미설정"과 "포트를 0으로 설정"을 구별할 수 없음
- 유효성 검사 위치가 애매함

---

## 해결: Functional Options 패턴

### 네 가지 구성 요소

**1. unexported `options` 구조체 — 설정값 보관**

```go
type options struct {
    port    *int          // 포인터: nil = 미설정, 0 = 랜덤 포트
    timeout time.Duration
    maxConn int
}
```

**2. `Option` 함수 타입 — options를 수정**

```go
type Option func(opts *options) error
```

**3. `WithXxx` 함수 — Option을 반환하고, 유효성 검사를 여기서**

```go
func WithPort(port int) Option {
    return func(opts *options) error {
        if port < 0 {
            return errors.New("port should be positive")
        }
        opts.port = &port
        return nil
    }
}
```

**4. 생성자에서 옵션을 순서대로 적용**

```go
func NewServer(addr string, opts ...Option) (*http.Server, error) {
    var o options
    for _, opt := range opts {
        if err := opt(&o); err != nil {
            return nil, err
        }
    }
    // 기본값 적용 ...
}
```

---

## 포인터 트릭

`int` 대신 `*int`를 사용해 세 가지 상태를 구별:

| 값 | 의미 |
|----|------|
| `nil` | 옵션 미제공 → 기본값 사용 |
| `&0`  | 명시적으로 0 설정 → 랜덤 포트 사용 |
| `&n`  | 명시적으로 n 설정 → n 사용 |

---

## 사용 예시

```go
NewServer("localhost")                                               // 모두 기본값
NewServer("localhost", WithPort(8080))                              // 포트만 지정
NewServer("localhost", WithPort(8080), WithTimeout(30*time.Second)) // 여러 옵션
NewServer("localhost", WithPort(-1))                                // 에러 반환
```

---

## 비교

| | 파라미터 폭발 | Config 구조체 | Functional Options |
|---|---|---|---|
| 옵션 추가 | 호출부 전체 수정 | 수정 없음 | 수정 없음 |
| 옵션 생략 | 기본값 직접 전달 | `Config{}` 전달 | 그냥 생략 |
| 유효성 검사 | 생성자 안 | 생성자 안 | 각 `WithXxx` 안 |
| nil vs 0 구별 | ❌ | ❌ | ✅ 포인터 사용 |
| Go 관용적 | ❌ | △ | ✅ |

---

## 핵심 원칙

- **필수값** → 일반 파라미터 (`addr string`)
- **선택값** → `Option` 함수 (`WithPort`, `WithTimeout`)
- **기본값 vs 미제공 구별** → 포인터(`*int`)로 `nil` 체크
- **유효성 검사** → 생성자가 아닌 각 `WithXxx` 안에서
