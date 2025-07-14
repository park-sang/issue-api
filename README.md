# 이슈 관리 API (Issue Management API)

## 프로젝트 개요
- Go 언어와 Gin 프레임워크를 사용한 이슈 관리 REST API입니다.
- 주요 기능: 이슈 생성, 조회, 수정, 상태 관리
- 포트: 8080에서 실행됩니다.

## 실행 방법

1. Go 설치 (버전 1.20 이상 권장)  
   https://go.dev/dl/

2. 프로젝트 클론 및 의존성 설치  
   git clone https://github.com/park-sang/issue-api.git  
   cd issue-api  
   go mod tidy  

3. 서버 실행  
   go run main.go  

   → 서버가 `localhost:8080` 에서 실행됩니다.

## API 명세 요약

### [POST] /issue - 이슈 생성
- 필수 필드: title
- 선택 필드: description, userId
- 담당자가 있으면 상태는 IN_PROGRESS, 없으면 PENDING

---

### [GET] /issues - 이슈 목록 조회
- 전체 이슈 반환
- status 쿼리 파라미터로 필터링 가능

---

### [GET] /issue/:id - 이슈 상세 조회

---

### [PATCH] /issue/:id - 이슈 수정
- 수정 가능 필드: title, description, status, userId
- 상태 전환 및 담당자 변경은 비즈니스 규칙에 따라 제한
- COMPLETED 또는 CANCELLED 상태는 수정 불가

---

## API 테스트 방법
- VSCode의 REST Client 확장 또는 Postman 추천
- 루트 디렉토리에 포함된 `test.http` 파일을 통해 테스트 가능

---

## 주요 구현 사항
- sync.Mutex를 이용한 동시성 제어
- 사용자 3명 하드코딩 (김개발, 이디자인, 박기획)
- 상태값 및 입력값 검증 철저
- 명확한 에러 메시지 및 적절한 HTTP 상태코드 제공
- Gin 프레임워크 기반 RESTful API 설계
