# go-hellogsm-ops

## 소개

go-hellogsm-ops는 www.hellogsm.kr (광주소프트웨어마이스터고 입학지원시스템)의 운영 보조 application을 모아둔 repo 입니다.

## relay-api

운영 관련 내용을 더모먼트팀 디스코드 웹훅 메시지로 보내고자 할 때 사용한다.

## relay-api 실행하기

```shell
go build./cmd/relay-api
./relay-api
```

## generate-dml

특정 조건에 부합하는 mock insert query 를 자동으로 추출해주는 프로그램이다.  
(1)로컬환경에서 개발, (2)어플리케이션 통합 테스트에 사용하면 유용하다.

## generate-dml 실행하기

```shell
go build./cmd/generate-dml
./generate-dml -graduate [] -screening [] -status []
```

### 파라미터 소개

```text
-graduate -screening -status

ex) 1차 배치 전 상태의 일반전형 지원자 100명, 사회통합전형 지원자 30명, 정원외특별전형 지원자 4명 만들어줘
./generate-dml -screening GEN100, SPE30,EXT4 -status FIRST
```

#### 졸업상태 [gradute]

- default - random
- 졸업예정 **[CANDIDATE]**
- 졸업자 **[GRADUATE]**
- 검정고시 **[GED]**

#### 전형 별 지원자 수 [screening] *required

- 일반전형, 사회통합전형, 정원외특별전형의 지원자 수 각각 **GEN**, **SPE**, **EXT**의 prefix 뒤에 입력합니다.
- 한번에 여러 전형의 원서를 생성하고 싶다면 “,” 컴마를 기준으로 여러 전형을 입력할 수 있습니다.
    - ex. -screening GEN10,SPE5,EXT2 → 일반전형 10명, 사회통합전형 5명, 정원외특별전형 2명
    - Extra는 국가보훈대상자 [EXTRA_ADMISSION], 특례입학대상자 [EXTRA_VETERANS] 중 하나가 랜덤으로 할당 됩니다.

#### 원서상태 [status] *required

- 1차 배치 전 base data **[FIRST]**
- 2차 배치 전 base data **[SECOND]**
- 최종 학과 배정 배치 전 base data **[FINAL_MAJOR]**
- 추가 모집 배치 전 base data **[RE_EVALUATE]**

