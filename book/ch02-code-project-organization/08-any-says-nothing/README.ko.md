# any는 아무 말도 안 한다 (#8)

## 한 줄 요약

`json.Marshal`처럼 진짜 어떤 타입이든 받아야 하는 경우에만 `any`를 사용하라. 그 외에는 `any`가 의미 있는 정보를 제공하지 않으며 컴파일 타임 안전성을 잃게 된다.

## `any`란?

`any`는 `interface{}`의 타입 별칭으로, 어떤 타입이든 담을 수 있다:

```go
var i any

i = 42
i = "foo"
i = struct{ s string }{s: "bar"}
i = someFunc
```

## 문제: `any` 남용

### Bad

```go
type Store struct{}

func (s *Store) Get(id string) (any, error) {
    return nil, nil
}

func (s *Store) Set(id string, v any) error {
    return nil
}
```

**문제점:**
- 컴파일 타임 타입 검사 없음
- 꺼낼 때 타입 단언 필요 → 런타임 패닉 위험
- 뭘 넣어야 하는지, 뭐가 나오는지 코드만 봐선 모름
- IDE 자동완성 안 됨

### Good

```go
func (s *Store) GetContract(id string) (Contract, error) {
    return Contract{}, nil
}

func (s *Store) SetContract(id string, contract Contract) error {
    return nil
}

func (s *Store) GetCustomer(id string) (Customer, error) {
    return Customer{}, nil
}

func (s *Store) SetCustomer(id string, customer Customer) error {
    return nil
}
```

**장점:**
- 컴파일 타임 타입 안전성
- 타입 단언 불필요
- 의도 명확
- IDE 지원 향상

## `any`를 써도 되는 경우

진짜 어떤 타입이든 다뤄야 할 때:

```go
json.Marshal(v any)       // JSON 마샬링
fmt.Println(a ...any)     // 포매팅
```

## 핵심 원칙

> **약간의 코드 중복이 `any` 남용보다 낫다.**

`any`는 코드 표현성 같은 다른 측면을 개선할 때 가끔 코드 중복이 더 나을 수 있다. 코드를 지나치게 일반화하는 것은 피해야 한다.
