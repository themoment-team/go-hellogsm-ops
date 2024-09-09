package main

import (
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"strings"
	_type "themoment-team/hellogsm-notice-server/generate-dml/type"
)

func main() {
	rowsParam := flag.Int("rows", 1, "Number of rows to generate")
	graduateStatusParam := flag.String("graduate", "RANDOM", "Status of grade")
	screeningParam := flag.String("screening", "RANDOM", "Status of screening")
	oneseoStatusParam := flag.String("status", "필수 입력 파라미터 입니다.", "Status of oneseo")

	flag.Parse()

	graduateStatus := _type.GraduateStatus(strings.ToUpper(*graduateStatusParam))
	screening := _type.Screening(strings.ToUpper(*screeningParam))
	oneseoStatus := _type.OneseoStatus(strings.ToUpper(*oneseoStatusParam))
	rows := *rowsParam

	err := validateParameter(graduateStatus, screening, oneseoStatus)
	if err != nil {
		fmt.Println(err)
		return
	}

	// GraduateStatus 가 RANDOM_GRADUATE_STATUS 라면 랜덤한 GraduateStatus 배열을 생성한 후 같은 인덱스의 row 들에 공통적으로 적용
	var graduateStatuses []_type.GraduateStatus
	if graduateStatus == _type.RANDOM_GRADUATE_STATUS {
		graduateStatuses = make([]_type.GraduateStatus, rows)
		for i := 0; i < rows; i++ {
			graduateStatuses[i] = []_type.GraduateStatus{
				_type.CANDIDATE,
				_type.GRADUATE,
				_type.GED,
			}[rand.Intn(3)]
		}
	} else {
		graduateStatuses = make([]_type.GraduateStatus, rows)
		for i := range graduateStatuses {
			graduateStatuses[i] = graduateStatus
		}
	}

	memberInsertQuery := GenerateMemberInsertQuery(rows)
	oneseoInsertQuery := GenerateOneseoInsertQuery(rows, screening, oneseoStatus)
	oneseoPrivacyDetailInsertQuery := GenerateOneseoPrivacyDetailInsertQuery(rows, graduateStatuses)

	fmt.Println(memberInsertQuery)
	fmt.Println(oneseoInsertQuery)
	fmt.Println(oneseoPrivacyDetailInsertQuery)
}

func validateParameter(graduateStatus _type.GraduateStatus, screening _type.Screening, oneseoStatus _type.OneseoStatus) error {
	if !graduateStatus.IsValidGraduateStatus() {
		return errors.New(fmt.Sprintf("잘못된 졸업상태가 입력되었습니다: %s", graduateStatus))
	}

	if !screening.IsValidScreening() {
		return errors.New(fmt.Sprintf("잘못된 전형이 입력되었습니다: %s", screening))
	}

	if !oneseoStatus.IsValidOneseoStatus() {
		return errors.New(fmt.Sprintf("잘못된 원서상태가 입력되었습니다: %s", oneseoStatus))
	}

	return nil
}
