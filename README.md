# go-hellogsm-ops

<br>

## 소개
go-hellogsm-ops는 www.hellogsm.kr (광주소프트웨어마이스터고 입학지원시스템)의 운영 보조 application을 모아둔 repo 입니다.

<br>
<br>

## notice-server
운영 서버에서 팀 다스코드로 알림을 전송하는 릴레이 서버입니다.

### 실행하기
```go
go build ./cmd/notice-server
./notice-server
```

<br>
<br>

## generate-dml
배치 application이나 개발 서버에서 테스트시 필요한 mock insert query DML을 생성하는 API입니다.

### generate-dml 실행하기

```go
go build ./cmd/generate-dml
./generate-dml -graduate [] -screening [] -status [] -rows []
```

### 파라미터 소개

```go
-graduate -screening -status -rows

ex) 1차 배치 전 상태의 일반전형 지원자 100명 만들어줘
./generate-dml -screening GENERAL -status FIRST -rows 100
```

#### 졸업상태 [gradute]
- default - random
- 졸업예정 **[CANDIDATE]**
- 졸업자 **[GRADUATE]**
- 검정고시 **[GED]**

#### 전형 [screening]
  - default - random
  - 일반전형 **[GENERAL]**
  - 사회통합전형 **[SPECIAL]**
  - 정원 외 특별 전형
      - 국가보훈대상자 **[EXTRA_ADMISSION]**
      - 특례입학대상자 **[EXTRA_VETERANS]**
#### 원서상태 [status] *required
- 1차 배치 전 base data **[FIRST]**
- 2차 배치 전 base data **[SECOND]**
- 최종 학과 배정 배치 전 base data **[FINAL_MAJOR]**
- 추가 모집 배치 전 base data **[RE_EVALUATE]**

#### row수 [rows]
- 생성할 데이터의 row 수
    - default - 1
