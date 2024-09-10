package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"themoment-team/hellogsm-notice-server/cmd/generate-dml/type"
)

func GenerateOneseoInsertQuery(rows int, initialScreening _type.Screening, oneseoStatus _type.OneseoStatus) string {
	var buffer bytes.Buffer
	majors := []string{"AI", "SW", "IOT"}
	var allScreenings = []_type.Screening{
		_type.GENERAL,
		_type.SPECIAL,
		_type.EXTRA_ADMISSION,
		_type.EXTRA_VETERANS,
	}

	buffer.WriteString("-- tb_oneseo" + "\n\n")

	for i := 1; i <= rows; i++ {

		// Screening이 RANDOM 이라면 GENERAL, SPECIAL, EXTRA_ADMISSION, EXTRA_VETERANS 중 하나로 랜덤 배정
		screening := initialScreening
		if screening == _type.RANDOM_SCREENING {
			screening = allScreenings[rand.Intn(len(allScreenings))]
		}

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
