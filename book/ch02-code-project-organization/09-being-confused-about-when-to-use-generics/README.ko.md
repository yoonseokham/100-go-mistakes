# 제네릭을 언제 써야 하는지 혼란 (#9)

## 한 줄 요약

제네릭(타입 파라미터)은 보일러플레이트 코드를 줄여준다. 하지만 구체적인 필요가 생길 때까지 미리 쓰지 마라. 그렇지 않으면 불필요한 추상화와 복잡성이 생긴다.

## 제네릭 이전의 문제

`map[string]int`의 키를 추출하는 함수를 `map[int]string`에도 쓰려면:

```go
func getKeys(m map[string]int) []string { ... }
// map[int]string에도 쓰고 싶다면?
```

제네릭 이전 선택지:
- 코드 중복
- 코드 생성
- 리플렉션

### `any` 사용 시 문제점

```go
func getKeys(m any) ([]any, error) {
    switch t := m.(type) {
    case map[string]int:
        // 추출 로직 복사
    case map[int]string:
        // 추출 로직 또 복사
    default:
        return nil, fmt.Errorf("unknown type: %T", t)
    }
}
```

- 각 케이스마다 보일러플레이트 중복
- 컴파일 타임 대신 런타임에 타입 검사
- `[]any` 반환 → 타입 안전성 상실

## 해결책: 타입 파라미터

```go
func getKeys[K comparable, V any](m map[K]V) []K {
    var keys []K
    for k := range m {
        keys = append(keys, k)
    }
    return keys
}
```

**타입 파라미터:**
- `K comparable` — 맵 키는 비교 가능해야 함 (`==`, `!=` 사용 가능)
- `V any` — 값은 어떤 타입이든 가능
- 컴파일 타임에 인스턴스화 → 타입 안전성 유지

```go
m := map[string]int{"one": 1, "two": 2}
keys := getKeys(m)  // Go가 타입 인수 추론
```

## 제약(Constraint)

### 커스텀 제약 예시

```go
type customConstraint interface {
   ~int | ~string
}

func getKeys[K customConstraint, V any](m map[K]V) []K { ... }
```

### `~int` vs `int`

- `int` — 정확히 `int` 타입만
- `~int` — 기반 타입이 `int`인 모든 타입 (`int`, `customInt` 등)

```go
type customInt int

func (i customInt) String() string {
   return strconv.Itoa(int(i))
}

// customInt는 ~int 제약을 만족 (기반 타입이 int이고 String() 구현)
```

## 자료구조에서의 제네릭

```go
type Node[T any] struct {
    Val  T
    next *Node[T]
}

func (n *Node[T]) Add(next *Node[T]) {
    n.next = next
}
```

**주의:** 타입 파라미터는 메서드가 아닌 함수에만 사용 가능:

```go
func (Foo) bar[T any](t T) {}  // ❌ 컴파일 에러
```

## 써야 할 때

1. **자료구조** — 이진 트리, 링크드 리스트, 힙 등

2. **어떤 타입의 slice/map/channel이든 처리**
```go
func merge[T any](ch1, ch2 <-chan T) <-chan T { ... }
```

3. **동작 팩터링**
```go
type sliceFn[T any] struct {
   s       []T
   compare func(T, T) bool
}
```

## 쓰면 안 될 때

### 타입 파라미터로 메서드만 호출하는 경우

```go
// ❌ 제네릭 불필요
func foo[T io.Writer](w T) {
   w.Write(getBytes())
}

// ✅ 인터페이스로 충분
func foo(w io.Writer) {
   w.Write(getBytes())
}
```

`T`의 타입 정보를 실제로 활용하지 않고 메서드만 호출한다면 인터페이스로 대체 가능.

### 코드를 더 복잡하게 만들 때

제네릭은 필수가 아니다. Go 개발자들은 10년 이상 제네릭 없이 살았다. 코드를 더 명확하게 만들지 못한다면 쓰지 마라.

## 인터페이스 남용과의 유사성

> 제네릭도 추상화의 한 형태 → 불필요한 추상화 = 복잡성 증가

## 핵심 원칙

- 보일러플레이트 코드가 생기려는 순간에 제네릭 도입
- 그 전에는 코드 중복이 더 나을 수 있음
- 타입 파라미터를 조기에 사용하지 마라
- 타입 자체를 보존/조작해야 할 때 제네릭, 메서드만 호출하면 인터페이스
