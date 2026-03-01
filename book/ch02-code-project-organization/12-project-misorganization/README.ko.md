# 프로젝트 구조 잘못 잡기 (#12)

## 한 줄 요약

성급하게 패키지를 나누지 마라. 패키지 이름은 "무엇을 담는지"가 아닌 "무엇을 제공하는지"로. export는 최소화하고, `internal/`로 구현 세부사항을 숨겨라. 구조(컨텍스트별/레이어별)를 하나 정해서 일관성 있게 유지하라.

---

## 프로젝트 구조: 두 가지 방식

### 컨텍스트별 구성 (도메인 기준)

비즈니스 개념(고객, 주문, 상품)을 기준으로 묶음:

```
shop/
  customer/
    handler.go    ← 고객 HTTP 핸들러
    service.go    ← 고객 비즈니스 로직
    repository.go ← 고객 DB 쿼리
  order/
    handler.go
    service.go
    repository.go
```

"고객 코드가 어디 있어?" → `customer/` 하나만 보면 됨.

### 레이어별 구성 (기술 계층 기준)

기술적 책임(HTTP, 비즈니스 로직, DB)을 기준으로 묶음:

```
shop/
  handler/        ← 모든 HTTP 핸들러
    customer.go
    order.go
  service/        ← 모든 비즈니스 로직
    customer.go
    order.go
  repository/     ← 모든 DB 쿼리
    customer.go
    order.go
```

"DB 쿼리 코드가 어디 있어?" → `repository/` 하나만 보면 됨.

**둘 다 정답. 핵심은 일관성 — 섞어 쓰는 게 제일 나쁨.**

---

## 세 가지 레이어

### Handler
- HTTP 요청 수신, 파싱, 응답 반환
- Service 호출
- 비즈니스 로직 없음 — 그냥 "중계자"

```go
func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
    json.NewDecoder(r.Body).Decode(&req)              // 요청 파싱
    customer, err := h.svc.CreateCustomer(req.Name)   // Service 위임
    json.NewEncoder(w).Encode(customer)               // 응답 반환
}
```

### Service
- 모든 비즈니스 규칙이 여기 있음
- 호출자가 HTTP인지 gRPC인지 모름
- DB가 MySQL인지 Postgres인지 모름
- Repository **인터페이스**만 앎

```go
func (s *CustomerService) CreateCustomer(name string) (Customer, error) {
    if name == "" {
        return Customer{}, errors.New("이름은 필수")  // 비즈니스 규칙
    }
    return s.repo.Save(Customer{Name: name})  // 어떤 DB인지 모름
}
```

### Repository
- DB 쿼리만 담당
- 비즈니스 로직 없음

```go
func (r *CustomerRepository) Save(c Customer) (Customer, error) {
    _, err := r.db.Exec("INSERT INTO customers (name) VALUES (?)", c.Name)
    return c, err
}
```

---

## 의존성 방향

의존성은 단방향으로만 흐름:

```
Handler → Service → Repository → DB
```

각 레이어는 아래 레이어를 **인터페이스**로만 앎 (구체 타입 모름):

```go
// Handler는 Service를 인터페이스로 앎
type customerService interface {
    CreateCustomer(name string) (Customer, error)
}

// Service는 Repository를 인터페이스로 앎
type CustomerRepo interface {
    Save(customer Customer) (Customer, error)
}
```

인터페이스는 **소비자 측**에서 정의 (6~7장 원칙).

### 의존성 연결은 main에서만

```go
func main() {
    db   := connectDB()
    repo := repository.NewCustomerRepository(db)  // DB 주입
    svc  := service.NewCustomerService(repo)      // repo 주입
    h    := handler.NewCustomerHandler(svc)       // svc 주입
}
```

이것이 **의존성 주입(DI)**. 5~7장에서 배운 "인터페이스를 파라미터로 받아라"가 실제로 적용되는 것.

---

## 왜 레이어를 나눠야 하는가

| 문제 (Bad) | 해결 (Good) |
|---|---|
| DB 없이 테스트 불가 | Service에 MockRepo 주입 |
| DB 교체 시 Handler 수정 필요 | Repository만 교체 |
| 비즈니스 규칙이 여기저기 흩어짐 | 모든 규칙이 Service에 |
| `init()`으로 테스트 어려움 | 명시적 생성자 사용 |

---

## 패키지 모범 사례

### 성급한 패키지 분리 금지

처음부터 완벽한 구조를 만들려 하지 마라. 코드가 커지면서 자연스럽게 분리.

### nano 패키지 금지

```
❌ 너무 잘게 나눔              ✅ 균형 잡힌 크기
  parse/parse.go               encoding/
  validate/validate.go           encode.go
  format/format.go               decode.go
                                 validate.go
```

### 패키지 이름: "무엇을 제공하는지" 기준

```
❌ utils, helpers, common, misc
✅ http, auth, encoding, storage
```

짧고, 간결하고, 표현적이고, 소문자 단어 하나.

### export 최소화

```go
// 기본: unexported
type internalHelper struct { ... }

// 필요할 때만 export
type PublicAPI struct { ... }
```

확신이 없으면 → export 하지 마라. 나중에 필요하면 그때 export.

### `internal/` 활용

```
project/
  internal/       ← 외부 모듈에서 import 불가
    service/
    repository/
  cmd/
    server/
      main.go
```

`internal/` 안의 코드는 외부 사용자를 신경 쓰지 않고 자유롭게 리팩토링 가능.

---

## 권장 서버 구조

```
project/
  go.mod
  internal/
    handler/
      customer.go
    service/
      customer.go
    repository/
      customer.go
  cmd/
    server/
      main.go
```
