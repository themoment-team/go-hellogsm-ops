package main

import (
	"bytes"
	"fmt"
	"math/rand"
)

func GenerateOneseoInsertQuery(rows int, generalCount int, specialCount int, extraCount int, oneseoStatus OneseoStatus) string {
	var buffer bytes.Buffer
	buffer.WriteString("-- tb_oneseo_insert" + "\n\n")

	majors := []string{"AI", "SW", "IOT"}

	for i := 1; i <= rows; i++ {

		// 전형 별 지원자 수를 파라미터로 받아 생성할 query 의 전형을 초기화
		screening := initScreening(&generalCount, &specialCount, &extraCount)

		// submitCodePrefix는 GENERAL - A, SPECIAL - B, EXTRA_ADMISSION, EXTRA_VETERANS - C
		var submitCodePrefix string
		switch screening {
		case GENERAL:
			submitCodePrefix = "A"
		case SPECIAL:
			submitCodePrefix = "B"
		case EXTRA_ADMISSION, EXTRA_VETERANS:
			submitCodePrefix = "C"
		}

		submitCode := fmt.Sprintf("%s-%d", submitCodePrefix, i)

		// first_desired_major, second_desired_major, third_desired_major 컬럼은 SW, IOT, AI 중 하나로 랜덤 배정
		selectedMajors := rand.Perm(len(majors))
		firstMajor := majors[selectedMajors[0]]
		secondMajor := majors[selectedMajors[1]]
		thirdMajor := majors[selectedMajors[2]]

		var appliedScreening, wantedScreening Screening

		appliedScreening = screening
		wantedScreening = screening

		// OneseoStatus가 FIRST라면 applied_screening 컬럼에는 NULL값 할당
		appliedScreeningStr := "NULL"
		if oneseoStatus != FIRST {
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

func initScreening(generalCount *int, specialCount *int, extraCount *int) Screening {
	if *generalCount > 0 {
		*generalCount--

		return GENERAL
	} else if *specialCount > 0 {
		*specialCount--

		return SPECIAL
	} else {
		*extraCount--

		randomValue := rand.Intn(2)
		if randomValue == 0 {
			return EXTRA_ADMISSION
		} else {
			return EXTRA_VETERANS
		}
	}
}
