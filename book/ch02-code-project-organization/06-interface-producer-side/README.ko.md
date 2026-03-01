# 생산자 측 인터페이스 (#6)

## 한 줄 요약

인터페이스를 소비자(클라이언트) 측에 두면 불필요한 추상화를 피할 수 있다.

## 핵심 원칙

> **"추상화는 발견되는 것이지, 만들어지는 것이 아니다."**

- ❌ 생산자가 모든 클라이언트에게 특정 추상화를 강제하는 건 잘못됨
- ✅ 클라이언트가 자신에게 필요한 추상화를 결정해야 함

## 문제: 생산자 측에서 큰 인터페이스 정의

### Bad: 생산자가 6개짜리 인터페이스 강제

```go
// store 패키지 (생산자)
type CustomerStorage interface {
    StoreCustomer(customer Customer) error
    GetCustomer(id string) (Customer, error)
    UpdateCustomer(customer Customer) error
    GetAllCustomers() ([]Customer, error)
    GetCustomersWithoutContract() ([]Customer, error)
    GetCustomersWithNegativeBalance() ([]Customer, error)
}
```

**문제점:**
1. 메서드 6개 = 약한 추상화
2. `GetAllCustomers()` 하나만 필요한 클라이언트도 6개 전부에 의존해야 함
3. 재사용성 낮음 — 큰 인터페이스를 만족시키기 어려움
4. 생산자가 추상화 수준을 강제 — 클라이언트 선택권 없음

### 클라이언트 입장

```go
// client 패키지 (소비자)
type Client struct {
    storage store.CustomerStorage  // ❌ 6개 메서드 전부 의존 강제
}

func (c *Client) ListCustomers() ([]store.Customer, error) {
    return c.storage.GetAllCustomers()  // 실제로는 1개만 사용!
}
```

테스트 시 사용하지도 않는 메서드 5개를 mock으로 구현해야 함.

## 해결책: 소비자 측 인터페이스

### Good: 생산자는 struct만, 소비자가 필요한 것만 정의

```go
// store 패키지 (생산자): struct만 제공
type Store struct { ... }

func (s *Store) StoreCustomer(customer Customer) error { ... }
func (s *Store) GetCustomer(id string) (Customer, error) { ... }
func (s *Store) UpdateCustomer(customer Customer) error { ... }
func (s *Store) GetAllCustomers() ([]Customer, error) { ... }
// ... 나머지 메서드들
```

```go
// client 패키지 (소비자): 필요한 1개만 정의
type customersGetter interface {
    GetAllCustomers() ([]store.Customer, error)
}

type Client struct {
    getter customersGetter  // 작고 집중된 인터페이스
}
```

**장점:**
1. 인터페이스 1개 = 강한 추상화
2. 실제로 필요한 것만 의존
3. mock이 쉬움 — 1개만 구현하면 됨
4. 재사용성 높음

### 다른 소비자, 다른 인터페이스

```go
// reporter 패키지 (다른 소비자)
type customerReader interface {
    GetCustomer(id string) (store.Customer, error)
    GetCustomersWithNegativeBalance() ([]store.Customer, error)
}
```

**같은 `Store`가 여러 다른 인터페이스를 동시에 만족!**

## 생산자 측 인터페이스가 괜찮은 경우

확실히 범용적으로 유용하다고 **알 때** (예상할 때가 아니라):

```go
// io 패키지 (생산자)
type Reader interface {
    Read(p []byte) (n int, err error)  // 메서드 1개
}

type Writer interface {
    Write(p []byte) (n int, err error)  // 메서드 1개
}
```

**왜 괜찮은가:**
1. 범용적으로 유용하다는 게 증명됨 (예상이 아님)
2. 최소화 — 메서드 1개씩
3. 높은 재사용성 — 구현하기 쉬움
4. 조합 가능 — `ReadWriter`로 합칠 수 있음

## 비교

| 생산자 측 인터페이스 | 소비자 측 인터페이스 |
|---------------------|---------------------|
| 모든 클라이언트에게 추상화 강제 | 클라이언트가 추상화 선택 |
| 특정 필요에 비해 너무 큼 | 정확한 필요에 맞춤 |
| 낮은 재사용성 | 높은 재사용성 |
| `io.Reader` (범용적으로 유용) | `customersGetter` (특정 필요) |

## 핵심 원칙

- **기본: 소비자 측** — 대부분의 경우 인터페이스는 클라이언트 측에
- **작을수록 강력** — "인터페이스가 클수록 추상화는 약해진다"
- **생산자는 concrete** — 구체 타입 제공, 추상화는 클라이언트가
- **예외: 범용 추상화** — `io.Reader`처럼 진짜 범용적일 때만
