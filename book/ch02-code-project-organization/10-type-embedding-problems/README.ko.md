# 타입 임베딩 문제 (#10)

## 한 줄 요약

타입 임베딩은 보일러플레이트 코드를 줄여준다. 하지만 숨겨야 할 필드나 동작이 외부에 노출되지 않도록 주의해야 한다.

## 타입 임베딩이란?

Go에서 struct 필드가 이름 없이 선언되면 **임베딩**이라고 한다:

```go
type Foo struct {
    Bar  // 임베딩 (이름 없음)
}

type Bar struct {
    Baz int
}
```

### 필드/메서드 승격(Promotion)

임베딩은 내부 타입의 필드와 메서드를 **승격**시킨다. `Bar`의 `Baz` 필드가 `Foo`로 승격:

```go
foo := Foo{}
foo.Baz = 42  // Bar.Baz에 직접 접근 가능
```

## 임베딩을 써야 하는 경우: 동작 합성

포워딩 메서드 보일러플레이트를 없앨 때:

### Bad: 포워딩 메서드 직접 작성

```go
type Logger struct {
    writeCloser io.WriteCloser
}

// Write, Close 쓰려면 위임 메서드를 일일이 작성해야 함
func (l Logger) Write(p []byte) (int, error) {
    return l.writeCloser.Write(p)  // 그냥 넘기는 것뿐
}

func (l Logger) Close() error {
    return l.writeCloser.Close()  // 그냥 넘기는 것뿐
}
```

### Good: 임베딩으로 자동 승격

```go
type Logger struct {
    io.WriteCloser  // Write(), Close() 자동으로 사용 가능
}

func main() {
    l := Logger{WriteCloser: os.Stdout}
    l.Write([]byte("foo"))  // 포워딩 메서드 없이 바로 사용
    l.Close()
}
```

**활용 예시 — 커스텀 동작 추가:**

```go
type Logger struct {
    io.WriteCloser
}

// Write만 오버라이드해서 타임스탬프 추가
func (l Logger) Write(p []byte) (int, error) {
    msg := append([]byte(time.Now().String()+" "), p...)
    return l.WriteCloser.Write(msg)  // 원본 Write 호출
}
// Close()는 그대로 임베딩에서 승격됨
```

## 임베딩을 쓰면 안 되는 경우

### 1. 단순 문법 편의를 위한 임베딩

```go
// ❌ Bad: foo.Baz라고 쓰려고 임베딩
type Foo struct {
    Bar
}
// foo.Bar.Baz 대신 foo.Baz를 쓰는 게 유일한 이유라면 → 이름 있는 필드 사용
```

### 2. 내부 구현 노출

숨겨야 할 데이터나 동작이 외부에 공개되면 안 됨:

#### Bad: sync.Mutex 임베딩

```go
type InMem struct {
    sync.Mutex  // ❌ Lock/Unlock이 public으로 노출됨!
    m map[string]int
}

inMem := New()
inMem.Lock()    // ❌ 외부에서 직접 호출 가능 → 오용 가능
inMem.Unlock()
```

뮤텍스는 내부 구현 세부사항인데 임베딩으로 외부에 노출됨.

#### Good: 이름 있는 필드 사용

```go
type InMem struct {
    mu sync.Mutex  // ✅ 소문자 → 외부 노출 안 됨
    m  map[string]int
}

// inMem.Lock()  → 컴파일 에러!
// mu는 내부에서만 사용
```

## 두 가지 핵심 제약

1. **단순 문법 편의로만 쓰지 마라** — `foo.Baz` vs `foo.Bar.Baz`가 유일한 이유라면 필드 사용
2. **숨겨야 할 것을 노출하지 마라** — 내부 잠금 동작 같은 건 외부에 보이면 안 됨

## 정리

| 상황 | 방법 |
|------|------|
| 포워딩 메서드 제거, 동작 합성 | 임베딩 (`io.WriteCloser`) |
| 내부 구현 숨기기 | 이름 있는 필드 (`mu sync.Mutex`) |
| 단순 문법 편의만 | 이름 있는 필드 사용 |

> **임베딩 = 동작(behavior) 합성 목적으로만. 구현 세부사항은 숨겨라.**
