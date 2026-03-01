# 인터페이스 오염 (#5)

## 한 줄 요약

추상화는 발견되는 것이지, 만들어지는 것이 아니다. 인터페이스는 필요할 때 만들어야 하며, 나중에 필요할 것 같다는 이유로 미리 만들면 안 된다.

## 인터페이스 오염이란?

불필요한 추상화로 코드를 가득 채워 이해하기 어렵게 만드는 것. Java, C# 같은 언어 습관을 가진 개발자들이 흔히 저지르는 실수다.

## Go 인터페이스 이해하기

### 암묵적 만족

Go 인터페이스는 **암묵적으로** 만족된다. `implements` 키워드가 필요 없다:

```go
// *os.File은 io.Reader를 자동으로 만족
file, _ := os.Open("file.txt")
var r io.Reader = file  // 그냥 동작함!
```

### 인터페이스 크기와 추상화 강도

> **"인터페이스가 클수록 추상화는 약해진다."** — Rob Pike

- 작은 인터페이스가 더 재사용하기 쉬움
- `io.Reader`, `io.Writer`가 강력한 이유: 더 이상 단순해질 수 없기 때문

## 인터페이스를 써야 하는 3가지 경우

### 1. 공통 동작 팩터링

여러 타입이 공통 동작을 구현할 때:

```go
// sort.Interface: 인덱스 기반 컬렉션이면 뭐든 정렬 가능
type Interface interface {
    Len() int
    Less(i, j int) bool
    Swap(i, j int)
}
```

### 2. 디커플링 (테스트 가능성)

구체 구현에서 코드를 분리할 때:

```go
// ❌ Bad: 구체 타입에 강하게 결합
type CustomerService struct {
    store *mysql.Store
}

// ✅ Good: 인터페이스에 의존
type customerStorer interface {
    StoreCustomer(Customer) error
}

type CustomerService struct {
    storer customerStorer  // MySQL, Postgres, MockStore 모두 주입 가능
}
```

### 3. 동작 제한

타입을 특정 동작으로 제한할 때:

```go
type IntConfig struct { ... }

func (c *IntConfig) Get() int { ... }
func (c *IntConfig) Set(value int) { ... }

// 읽기 전용으로 제한
type intConfigGetter interface {
    Get() int
}

type Foo struct {
    threshold intConfigGetter  // Set() 호출 불가
}
```

## 인터페이스 오염 문제

### Bad: 불필요한 인터페이스

```go
type UserService interface {
    CreateUser(name string) error
    GetUser(id int) (string, error)
}

type UserServiceImpl struct {
    db *sql.DB
}

// 구현체가 하나뿐인데 인터페이스를 만든 이유가 없음!
```

**왜 나쁜가:**
- 불필요한 간접 레이어
- 가치 없는 추상화
- 코드 탐색 어려움 (interface → impl → 실제 코드)
- 인터페이스가 존재할 이유 없음

### 핵심 원칙

> **"추상화는 발견되는 것이지, 만들어지는 것이 아니다."**

- 즉각적인 이유 없이 추상화 만들지 마라
- 인터페이스가 아닌 구체 타입으로 설계하라
- **필요할 때** 인터페이스를 만들어라, **필요할 것 같을 때**가 아니라

## Good Example: 발견된 추상화

### Step 1: 구체 타입으로 시작

```go
type UserService struct {
    db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
    return &UserService{db: db}
}
```

### Step 2: 필요할 때 소비자 패키지에서 인터페이스 정의

```go
// handler 패키지 (소비자)
type UserCreator interface {
    CreateUser(name string) error
}

type UserHandler struct {
    creator UserCreator
}
```

### Step 3: 자동으로 만족됨

```go
userSvc := userservice.NewUserService(db)
handler := handler.NewUserHandler(userSvc)  // 자동으로 인터페이스 만족
```

## 정리

| 상황 | 인터페이스 사용? |
|------|----------------|
| 구현체 하나, 이유 없음 | ❌ |
| "나중에 필요할 것 같아서" | ❌ |
| 테스트를 위한 디커플링 | ✅ |
| 동작(접근 권한) 제한 | ✅ |
| 공통 동작 여러 타입에 팩터링 | ✅ |
