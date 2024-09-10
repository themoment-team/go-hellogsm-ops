package main

import (
	"bytes"
	"fmt"
	"themoment-team/hellogsm-notice-server/cmd/generate-dml/type"
)

func GenerateEntranceTestResultInsertQuery(rows int, oneseoStatus _type.OneseoStatus, totalSubjectsScores, totalNonSubjectsScores []float64) string {
	var buffer bytes.Buffer

	buffer.WriteString("-- tb_entrance_test_result" + "\n\n")

	for i := 1; i <= rows; i++ {
		var aptitudeEvaluationScore, interviewScore *float64
		var documentEvaluationScore float64
		var firstTestPassYn, secondTestPassYn *string

		// tb_entrance_test_factors_detail DML 생성시 배열에 저장해둔 교과성적, 비교과성적을 더해 서류전형 총점을 계산
		documentEvaluationScore = totalSubjectsScores[i-1] + totalNonSubjectsScores[i-1]

		switch oneseoStatus {
		case _type.FIRST:
			aptitudeEvaluationScore = nil
			interviewScore = nil
			firstTestPassYn = nil
			secondTestPassYn = nil

		case _type.SECOND:
			aptitudeEvaluationScoreValue := randomFloat(0, 100)
			interviewScoreValue := randomFloat(0, 100)

			aptitudeEvaluationScore = &aptitudeEvaluationScoreValue
			interviewScore = &interviewScoreValue

			firstTestPassYn = stringPointer("YES")
			secondTestPassYn = nil

		case _type.FINAL_MAJOR:
			aptitudeEvaluationScoreValue := randomFloat(0, 100)
			interviewScoreValue := randomFloat(0, 100)

			aptitudeEvaluationScore = &aptitudeEvaluationScoreValue
			interviewScore = &interviewScoreValue

			firstTestPassYn = stringPointer("YES")
			secondTestPassYn = stringPointer("YES")

		case _type.RE_EVALUATE:
			aptitudeEvaluationScoreValue := randomFloat(0, 100)
			interviewScoreValue := randomFloat(0, 100)

			aptitudeEvaluationScore = &aptitudeEvaluationScoreValue
			interviewScore = &interviewScoreValue

			firstTestPassYn = stringPointer("YES")
			secondTestPassYn = stringPointer("NO")
		}

		query := fmt.Sprintf(
			"INSERT INTO tb_entrance_test_result (aptitude_evaluation_score, document_evaluation_score, interview_score, entrance_test_factors_detail_id, oneseo_id, first_test_pass_yn, second_test_pass_yn) "+
				"VALUES (%s, %.2f, %s, %d, %d, %s, %s);",
			formatNullable(aptitudeEvaluationScore),
			documentEvaluationScore,
			formatNullable(interviewScore),
			i,
			i,
			formatNullableString(firstTestPassYn),
			formatNullableString(secondTestPassYn),
		)

		buffer.WriteString(query + "\n")
	}

	return buffer.String()
}

func stringPointer(value string) *string {
	return &value
}

func formatNullableString(value *string) string {
	if value == nil {
		return "NULL"
	}
	return fmt.Sprintf("'%s'", *value)
}
