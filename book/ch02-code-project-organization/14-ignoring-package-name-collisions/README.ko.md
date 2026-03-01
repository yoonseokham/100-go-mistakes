# 패키지 이름 충돌 무시하기 (#14)

## 한 줄 요약

변수 이름과 패키지 이름이 충돌하면 패키지를 못 쓰게 된다. 고유한 변수 이름을 짓는 게 최선이고, import alias는 최후의 수단으로만.

---

## 문제

변수 이름이 패키지 이름과 같으면 그 패키지를 더 이상 쓸 수 없음:

```go
import "redis"

func process() {
    redis := redis.NewClient()  // ❌ 이 순간부터 redis는 변수, 패키지 아님
    redis.Ping()                // 변수 redis를 가리킴

    // 이 아래에서 redis.NewClient() 다시 호출 불가
    // redis는 이제 변수 이름이니까
}
```

---

## 해결책 (우선순위 순)

### 1순위: 변수 이름 바꾸기 (가장 좋음)

가장 단순하고 가독성 좋은 방법:

```go
redisClient := redis.NewClient()   // 명확하고 구체적
```

### 2순위: 더 구체적인 이름 짓기

```go
customerCache := redis.NewClient()  // 용도까지 드러남
```

### 3순위: import alias (최후의 수단)

```go
import redisapi "redis"

redis := redisapi.NewClient()
```

**왜 최후의 수단?**
- 읽는 사람이 `redisapi`를 보고 "이게 뭐지?" 혼란
- 파일마다 `redis`, `redisapi` 섞이면 일관성 깨짐
- 실질적 이점 없이 인지 부하만 추가

---

## import alias가 괜찮은 경우

두 패키지 이름이 진짜 똑같을 때 — 다른 방법이 없음:

```go
import (
    "crypto/rand"
    mrand "math/rand"  // 둘 다 "rand"라서 어쩔 수 없음
)
```

---

## 흔한 충돌 사례

```go
// ❌ context 패키지와 충돌
context := req.Context()
// ✅
ctx := req.Context()

// ❌ http 패키지와 충돌
http := doRequest()
// ✅
resp := doRequest()

// ❌ json 패키지와 충돌
json := marshal(data)
// ✅
data := marshal(data)

// ❌ sql 패키지와 충돌
sql := buildQuery()
// ✅
query := buildQuery()
```

---

## 핵심 원칙

대부분의 충돌은 변수 이름을 잘 지으면 피할 수 있다. import alias는 두 임포트 패키지 이름이 같아서 방법이 없을 때만 사용.

| 우선순위 | 방법 | 사용 시점 |
|---------|------|----------|
| 1순위 | 변수 이름 바꾸기 | 항상 먼저 시도 |
| 2순위 | 더 구체적인 이름 | 용도가 드러나면 더 좋을 때 |
| 3순위 | import alias | 두 패키지 이름이 같을 때만 |
