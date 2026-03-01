# init 함수 남용 (#3)

## 한 줄 요약

init 함수는 에러 처리가 제한적이고 테스트를 복잡하게 만든다. 대부분의 초기화는 명시적인 함수로 처리해야 한다.

## init 함수란?

- 인수도 반환값도 없는 함수 (`func()` 시그니처)
- 패키지가 초기화될 때 자동으로 실행
- `main()` 함수보다 먼저 실행

### 실행 순서

패키지 초기화 시:
1. 모든 상수/변수 선언 평가
2. 모든 `init()` 함수 실행 (선언 순서대로)
3. `main()` 실행

여러 `init()`이 있을 경우: 소스 파일 내 순서대로, 여러 파일이면 파일명 알파벳 순서로 실행.

## init 함수의 문제점

### 1. 에러 처리 불가

init 함수는 에러를 반환할 수 없다. 선택지는 panic 또는 에러 무시뿐:

```go
func init() {
    var err error
    db, err = sql.Open("postgres", "connection-string")
    if err != nil {
        panic(err)  // panic 아니면 무시, 선택지가 없음
    }
}
```

### 2. 테스트 복잡성

- 테스트 전에 자동으로 실행됨
- mock이나 초기화 제어 불가
- 불필요한 외부 의존성도 항상 셋업해야 함

```go
// Bad: 테스트가 init() 부작용에 의존
func TestStore(t *testing.T) {
    redis.Store("key", "value")  // init()이 성공적으로 실행됐다고 가정
}
```

### 3. 전역 상태 강제

```go
var cache map[string]string  // 전역 변수

func init() {
    cache = make(map[string]string)  // 전역 상태 설정
}
```

전역 상태는 예상치 못한 동작, 어려운 단위 테스트, 이해하기 어려운 코드를 만든다.

## 해결책: 명시적 초기화

### Good: 생성자 패턴

```go
type Cache struct {
    data map[string]string
}

func NewCache() (*Cache, error) {
    // 실패 시 에러 반환 가능
    return &Cache{
        data: make(map[string]string),
    }, nil
}
```

**장점:**
- 에러 처리 가능
- 여러 인스턴스로 테스트 쉬움
- 전역 상태 없음
- 초기화 시점 명확
- 의존성 주입 가능

### 테스트가 쉬워짐

```go
func TestCache(t *testing.T) {
    // 격리된 새 인스턴스 생성
    cache, err := NewCache()
    if err != nil {
        t.Fatalf("NewCache failed: %v", err)
    }
    cache.Store("key", "value")
}
```

## init을 써도 되는 경우

`init()`은 이런 경우에만:

1. **정적 설정** — 에러 가능성 없음
2. **레지스트리 등록** — 예: 데이터베이스 드라이버
3. **상수 데이터 셋업** — 결정론적, 부작용 없음

```go
var supportedFormats = map[string]bool{}

func init() {
    // OK: 단순하고 결정론적, 에러 없음
    supportedFormats["json"] = true
    supportedFormats["xml"] = true
}
```

## 핵심 원칙

- ❌ 실패할 수 있는 작업에 `init()` 사용 금지
- ❌ 무거운 작업에 `init()` 사용 금지
- ❌ 의존성 주입이 필요한 경우 `init()` 사용 금지
- ✅ 에러를 반환하는 명시적 생성자 사용
- ✅ 테스트 가능성을 위한 의존성 주입
- ✅ 초기화를 명시적이고 제어 가능하게
