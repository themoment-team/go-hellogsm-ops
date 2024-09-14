# go-hellogsm-ops

## 소개
go-hellogsm-ops는 www.hellogsm.kr (광주소프트웨어마이스터고 입학지원시스템)의 운영 보조 application을 모아둔 repo 입니다.

## 개발자 가이드

## notice-server
운영 서버에서 팀 다스코드로 알림을 전송하는 릴레이 서버입니다.

### notice-server 실행하기
```go
go build ./cmd/notice-server
./notice-server
```

## generate-dml
배치 application이나 개발 서버에서 테스트시 필요한 mock insert query DML을 생성하는 API입니다.

### generate-dml 실행하기

```go
go build ./cmd/generate-dml
./generate-dml -graduate [] -screening [] -status []
```

### 파라미터 소개

```go
-graduate -screening -status

ex) 1차 배치 전 상태의 일반전형 지원자 100명, 사회통합전형 지원자 30명, 정원외특별전형 지원자 4명 만들어줘
./generate-dml -screening 100,30,4 -status FIRST
```

#### 졸업상태 [gradute]
- default - random
- 졸업예정 **[CANDIDATE]**
- 졸업자 **[GRADUATE]**
- 검정고시 **[GED]**

#### 전형 별 지원자 수 [screening] *required
- 일반전형, 사회통합전형, 정원외특별전형의 지원자 수를 “,” 컴마를 기준으로 입력 받습니다. (10(General), 5(Special), 2(Extra))
  - ex. -screening 10,5,2 → 일반전형 10명, 사회통합전형 5명, 정원외특별전형 2명 
  - Extra는 국가보훈대상자 [EXTRA_ADMISSION], 특례입학대상자 [EXTRA_VETERANS] 중 하나가 랜덤으로 할당 됩니다.

#### 원서상태 [status] *required
- 1차 배치 전 base data **[FIRST]**
- 2차 배치 전 base data **[SECOND]**
- 최종 학과 배정 배치 전 base data **[FINAL_MAJOR]**
- 추가 모집 배치 전 base data **[RE_EVALUATE]**
