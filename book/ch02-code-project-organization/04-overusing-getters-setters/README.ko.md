# getter/setter 남용 (#4)

## 한 줄 요약

getter/setter를 강제로 사용하는 것은 Go답지 않다. 효율성과 관용적 코드 사이의 균형을 실용적으로 찾아야 한다.

## 데이터 캡슐화란?

객체의 값이나 상태를 숨기는 것. getter/setter는 unexported 필드 위에 exported 메서드를 제공해 캡슐화를 구현하는 수단이다.

### 다른 언어의 습관 (Java)

```java
class Person {
    private String name;

    public String getName() { return name; }
    public void setName(String name) { this.name = name; }
}
```

## Go에서의 getter/setter

### 자동 지원 없음

Go에는 getter/setter에 대한 **자동 지원이 없다**.

### 의무도 아니고 관용적이지도 않음

struct 필드 접근에 getter/setter를 쓰는 건 **의무도 관용적이지도 않다**.

### Bad: 불필요한 getter/setter

```go
type Person struct {
    name string
}

// 불필요: 그냥 필드 접근을 감싼 것뿐
func (p *Person) GetName() string {
    return p.name
}

func (p *Person) SetName(name string) {
    p.name = name
}
```

### Good: 직접 필드 접근

```go
type Person struct {
    Name string  // exported 필드, 직접 접근
}

person := Person{Name: "John"}
fmt.Println(person.Name)  // 직접 접근
person.Name = "Jane"      // 직접 수정
```

## Go 네이밍 컨벤션

getter/setter가 필요하다면 Go 컨벤션을 따를 것:

### Getter: "Get" 접두사 생략

```go
type Account struct {
    balance float64
}

// ✅ Good: Balance (GetBalance 아님)
func (a *Account) Balance() float64 {
    return a.balance
}

// ❌ Bad: "Get" 접두사 사용
func (a *Account) GetBalance() float64 {
    return a.balance
}
```

### Setter: "Set" 접두사 사용

```go
// ✅ Good: SetBalance
func (a *Account) SetBalance(balance float64) {
    a.balance = balance
}
```

## getter/setter를 써야 하는 경우

실제 가치를 제공할 때만 사용:

- 유효성 검사 로직 추가
- 접근 시 데이터 변환
- 계산된 값 반환
- 미래 변경에 대한 하위 호환성 보장
- 단순 필드 접근 이상의 동작 추가

## 정리

| 방법 | 사용 시점 |
|------|-----------|
| 직접 필드 접근 | 기본 선택 — 단순하고 관용적 |
| getter/setter | 실제 가치를 제공할 때만 |

**핵심:** 다른 언어의 패턴을 맹목적으로 따르지 마라. Go의 단순성 철학에 맞게 실용적으로 선택하라.
