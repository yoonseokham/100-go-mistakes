# Go 100 Mistakes 예제 실습

이 프로젝트는 "100 Go Mistakes and How to Avoid Them" 책의 예제 코드를 Bazel로 빌드/테스트합니다.

## 설치

Bazel 설치:
```bash
# macOS
brew install bazelisk
```

## 사용법

### 빌드
```bash
# 모든 타겟 빌드
bazel build //...

# 특정 예제 빌드
bazel build //examples/mistake01:example
```

### 실행
```bash
# 바이너리 실행
bazel run //examples/mistake01:example
```

### 테스트
```bash
# 모든 테스트 실행
bazel test //...

# 특정 테스트 실행
bazel test //examples/mistake01:mistake01_test
```

### Gazelle (자동 BUILD 파일 생성)
```bash
# BUILD.bazel 파일 자동 생성/업데이트
bazel run //:gazelle

# go.mod 기반으로 외부 의존성 업데이트
bazel run //:gazelle-update-repos
```

## 프로젝트 구조

```
.
├── WORKSPACE              # Bazel 워크스페이스 설정
├── BUILD.bazel           # 루트 BUILD 파일
├── go.mod                # Go 모듈 정의
├── .bazelrc              # Bazel 설정
└── examples/
    └── mistake01/
        ├── BUILD.bazel   # 예제별 빌드 규칙
        ├── lib.go        # 라이브러리 코드
        ├── lib_test.go   # 테스트 코드
        └── main.go       # 실행 가능한 예제
```

## Bazel 빌드 룰

- `go_library`: Go 패키지/라이브러리
- `go_binary`: 실행 가능한 바이너리
- `go_test`: Go 테스트
