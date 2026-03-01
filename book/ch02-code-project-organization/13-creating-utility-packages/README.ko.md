# 유틸리티 패키지 만들지 마라 (#13)

## 한 줄 요약

`common`, `util`, `shared` 같은 패키지를 만들면 읽는 사람에게 아무 정보도 주지 않는다. 의미 있고 구체적인 패키지 이름으로 리팩토링하라.

---

## 문제

### Bad: 아무 의미 없는 패키지 이름

```
project/
  utils/
    string.go     ← 문자열 헬퍼
    time.go       ← 시간 헬퍼
    convert.go    ← 변환 헬퍼
  common/
    errors.go
    constants.go
  shared/
    types.go
```

**사용할 때:**
```go
utils.NewStringSet("a", "b")     // "utils"가 뭔데?
utils.FormatTime(t)              // utils가 다 하네?
common.ErrNotFound               // common = 쓰레기통
```

**문제점:**
- 패키지 이름이 아무 의미도 전달 안 함 — "utils"에는 뭐든 들어갈 수 있음
- 관련 없는 코드가 한 패키지에 섞임
- 끝없이 커지는 쓰레기통이 됨
- 호출자가 패키지 목적을 한눈에 파악 불가

---

## 해결책

### Good: "무엇을 제공하는지" 기준으로 이름 짓기

```go
package stringset

type Set map[string]struct{}

func New(values ...string) Set { ... }
func (s Set) Sort() []string { ... }
```

**사용할 때:**
```go
stringset.New("a", "b", "c")   // 뭘 하는지 바로 보임
```

### 더 많은 예시

```
❌ utils.FormatTime()      → ✅ timeutil.Format()
❌ utils.ParseJSON()       → ✅ jsonparser.Parse()
❌ common.ErrNotFound      → ✅ customerr.ErrNotFound
❌ helpers.Validate()      → ✅ validation.Run()
```

---

## 표준 라이브러리도 이렇게 함

| 패키지 | 제공하는 것 |
|--------|-----------|
| `strings` | 문자열 조작 |
| `strconv` | 문자열 변환 |
| `net/http` | HTTP 클라이언트/서버 |
| `encoding/json` | JSON 인코딩/디코딩 |
| `net/http/httputil` | HTTP 유틸리티 (그냥 "utils" 아님) |

`httputil`조차도 범위가 한정됨 — HTTP 유틸리티이지, 범용 유틸리티가 아님.

---

## 핵심 원칙

패키지 이름은 **"무엇을 담는지"**가 아닌 **"무엇을 제공하는지"** 기준:

```
❌ 담고 있는 것 기준       ✅ 제공하는 것 기준
  utils                      stringset
  helpers                    httputil
  common                     auth
  shared                     encoding
  base                       storage
  misc                       validation
```

패키지 이름은 **짧고, 간결하고, 표현적이고, 소문자 단어 하나**.
