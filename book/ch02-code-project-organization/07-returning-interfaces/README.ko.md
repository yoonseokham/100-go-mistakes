# 인터페이스 반환 (#7)

## 한 줄 요약

유연성을 제한하지 않으려면 함수는 대부분의 경우 인터페이스가 아닌 구체 타입을 반환해야 한다. 반대로 함수는 가능하면 인터페이스를 파라미터로 받아야 한다.

## 핵심 원칙

> **"구체 타입을 반환하고, 인터페이스를 파라미터로 받아라."**

- **반환값**: 구체 타입 (struct, 구체 구현)
- **파라미터**: 인터페이스 (추상화)
- **추상화**: 생산자가 강제하는 게 아니라 클라이언트가 발견해야 함

## 문제: 인터페이스 반환

### Bad: 함수가 인터페이스 반환

```go
// store 패키지 (생산자)
type CustomerStorage interface {
    StoreCustomer(customer Customer) error
    GetCustomer(id string) (Customer, error)
    UpdateCustomer(customer Customer) error
    GetAllCustomers() ([]Customer, error)
}

// ❌ 인터페이스 반환
func NewInMemoryStore() CustomerStorage {
    return &InMemoryStore{
        customers: make(map[string]Customer),
    }
}
```

**문제점:**
1. `InMemoryStore`에 `Debug()` 메서드가 있어도 클라이언트가 쓸 수 없음
2. 인터페이스에 정의된 메서드만 접근 가능 → 기능 제한
3. 새 기능 추가 시 인터페이스 변경 필요 → 모든 클라이언트 영향
4. 패키지 간 결합도 증가

### 클라이언트의 한계

```go
type Client struct {
    storage store.CustomerStorage  // ❌ 이 인터페이스에 묶임
}

func NewClient() *Client {
    return &Client{
        storage: store.NewInMemoryStore(),
    }
}

// InMemoryStore에 Debug() 메서드가 있어도 사용 불가!
```

## 해결책: 구체 타입 반환

### Good: 함수가 구체 타입 반환

```go
// ✅ 구체 타입 반환
func NewInMemoryStore() *InMemoryStore {
    return &InMemoryStore{
        customers: make(map[string]Customer),
    }
}

// Debug() 같은 추가 메서드도 접근 가능
func (s *InMemoryStore) Debug() string { ... }
```

### 클라이언트가 자신의 인터페이스 정의

```go
// client 패키지
type customerGetter interface {
    GetCustomer(id string) (store.Customer, error)  // 필요한 것만
}

type Client struct {
    getter customerGetter
}

// ✅ 인터페이스를 파라미터로 받음
func NewClient(getter customerGetter) *Client {
    return &Client{getter: getter}
}
```

### 사용 예시

```go
// main.go
concreteStore := store.NewInMemoryStore()  // 구체 타입

client := client.NewClient(concreteStore)  // 인터페이스 자동 만족

debugInfo := concreteStore.Debug()  // 추가 메서드도 접근 가능
```

## 패턴 요약

```go
// ❌ BAD: 인터페이스 반환, 구체 타입 파라미터
func NewService(db *sql.DB) ServiceInterface { ... }

// ✅ GOOD: 구체 타입 반환, 인터페이스 파라미터
func NewService(db Database) *Service { ... }
```

## 인터페이스 반환이 괜찮은 경우

추상화가 범용적으로 유용하다고 **확실히 알 때** (예상이 아닌 확신):

```go
os.Open()  // io.ReadCloser 반환 — 이건 OK
```

**왜 괜찮은가:**
1. 범용적으로 유용하다는 게 입증됨
2. 최소화 — 메서드가 적음
3. stdlib 컨벤션과 일치

## 비교

| 인터페이스 반환 | 구체 타입 반환 |
|----------------|---------------|
| 클라이언트에게 추상화 강제 | 클라이언트가 추상화 선택 |
| 메서드 접근 제한 | 모든 메서드 접근 가능 |
| 확장하기 어려움 | 새 메서드 추가 쉬움 |
| 패키지 결합도 높음 | 느슨한 결합 |
| 생산자가 추상화 결정 | 소비자가 추상화 결정 |

## 핵심 원칙

- **기본: 구체 타입 반환** — 클라이언트에게 최대 유연성 제공
- **인터페이스는 파라미터로** — 어떤 호환 타입이든 전달 가능
- **클라이언트가 인터페이스 정의** — 자신의 필요를 가장 잘 아는 건 클라이언트
- **예외: 입증된 추상화** — stdlib 패턴처럼 확실히 범용적일 때만
