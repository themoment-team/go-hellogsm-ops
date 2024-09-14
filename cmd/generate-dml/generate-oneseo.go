package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"themoment-team/hellogsm-notice-server/cmd/generate-dml/type"
)

func GenerateOneseoInsertQuery(rows int, screeningCountArr []int, oneseoStatus _type.OneseoStatus) string {
	var buffer bytes.Buffer
	buffer.WriteString("-- tb_oneseo_insert" + "\n\n")

	majors := []string{"AI", "SW", "IOT"}

	generalScreeningCount := screeningCountArr[0]
	specialScreeningCount := screeningCountArr[1]
	extraScreeningCount := screeningCountArr[2]

	for i := 1; i <= rows; i++ {

		// 전형 별 지원자 수를 파라미터로 받아 생성할 query 의 전형을 초기화
		screening := initScreening(&generalScreeningCount, &specialScreeningCount, &extraScreeningCount)

		// submitCodePrefix는 GENERAL - A, SPECIAL - B, EXTRA_ADMISSION, EXTRA_VETERANS - C
		var submitCodePrefix string
		switch screening {
		case _type.GENERAL:
			submitCodePrefix = "A"
		case _type.SPECIAL:
			submitCodePrefix = "B"
		case _type.EXTRA_ADMISSION, _type.EXTRA_VETERANS:
			submitCodePrefix = "C"
		}

		submitCode := fmt.Sprintf("%s-%d", submitCodePrefix, i)

		// first_desired_major, second_desired_major, third_desired_major 컬럼은 SW, IOT, AI 중 하나로 랜덤 배정
		selectedMajors := rand.Perm(len(majors))
		firstMajor := majors[selectedMajors[0]]
		secondMajor := majors[selectedMajors[1]]
		thirdMajor := majors[selectedMajors[2]]

		var appliedScreening, wantedScreening _type.Screening

		appliedScreening = screening
		wantedScreening = screening

		// OneseoStatus가 FIRST라면 applied_screening 컬럼에는 NULL값 할당
		appliedScreeningStr := "NULL"
		if oneseoStatus != _type.FIRST {
			appliedScreeningStr = fmt.Sprintf("'%s'", appliedScreening)
		}

		query := fmt.Sprintf(
			"INSERT INTO tb_oneseo (member_id, applied_screening, oneseo_submit_code, first_desired_major, real_oneseo_arrived_yn, second_desired_major, third_desired_major, wanted_screening) "+
				"VALUES (%d, %s, '%s', '%s', 'YES', '%s', '%s', '%s');",
			i, appliedScreeningStr, submitCode, firstMajor, secondMajor, thirdMajor, wantedScreening,
		)

		buffer.WriteString(query + "\n")
	}

	return buffer.String()
}

func initScreening(generalScreeningCount *int, specialScreeningCount *int, extraScreeningCount *int) _type.Screening {
	if *generalScreeningCount > 0 {
		*generalScreeningCount--

		return _type.GENERAL
	} else if *specialScreeningCount > 0 {
		*specialScreeningCount--

		return _type.SPECIAL
	} else if *extraScreeningCount > 0 {
		*extraScreeningCount--

		randomValue := rand.Intn(2)
		if randomValue == 0 {
			return _type.EXTRA_ADMISSION
		} else {
			return _type.EXTRA_VETERANS
		}
	} else {
		panic("발생할 수 없는 상황입니다. 전형 별 지원자 수의 총 합과 rows의 수는 동일합니다.")
	}
}
