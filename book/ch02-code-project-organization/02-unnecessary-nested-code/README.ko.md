# 불필요한 중첩 코드 (#2)

## 문제

불필요한 `else` 블록과 깊은 들여쓰기는 코드를 읽기 어렵게 만든다. 중첩 레벨이 늘어날수록 인지 부하가 증가한다.

### Bad: 깊은 중첩

```go
func join1(s1, s2 string, max int) (string, error) {
    if s1 == "" {
        return "", errors.New("s1 is empty")
    } else {
        if s2 == "" {
            return "", errors.New("s2 is empty")
        } else {
            concat, err := concatenate(s1, s2)
            if err != nil {
                return "", err
            } else {
                if len(concat) > max {
                    return concat[:max], nil
                } else {
                    return concat, nil
                }
            }
        }
    }
}
```

**문제점:**
- 4단계 중첩
- `return` 이후 불필요한 `else`
- 핵심 로직(happy path)이 깊이 묻혀 있음

## 해결책: Early Return과 플랫 구조

### Good: Early Return, else 없음

```go
func join2(s1, s2 string, max int) (string, error) {
    if s1 == "" {
        return "", errors.New("s1 is empty")
    }
    if s2 == "" {
        return "", errors.New("s2 is empty")
    }
    concat, err := concatenate(s1, s2)
    if err != nil {
        return "", err
    }
    if len(concat) > max {
        return concat[:max], nil
    }
    return concat, nil
}
```

**장점:**
- 중첩 없음
- `else` 없음
- happy path가 왼쪽 정렬로 명확히 보임

## 핵심 원칙

1. **Early return 사용**: 에러 조건은 즉시 리턴
2. **return 후 else 제거**: `if` 블록에서 리턴하면 `else`는 필요 없음
3. **Happy path는 왼쪽**: 핵심 로직은 왼쪽 정렬, 맨 아래에
4. **Guard clause 먼저**: 입력 유효성 검사는 맨 위에서 처리

## 추가 패턴

### 루프에서 continue 활용

```go
for _, user := range users {
    if user == nil {
        continue
    }
    if !isActive(user) {
        continue
    }
    // 유효한 user 처리
}
```

### 중첩 검증 대신 early return

```go
// ❌ Bad: 중첩 검증
if valid {
    if authorized {
        if hasPermission {
            // 로직
        }
    }
}

// ✅ Good: early return
if !valid {
    return errors.New("invalid")
}
if !authorized {
    return errors.New("unauthorized")
}
if !hasPermission {
    return errors.New("no permission")
}
// 로직
```
