package main

import (
	"bytes"
	"fmt"
)

func GenerateEntranceTestFactorsDetailInsertQuery(rows int, graduateStatuses []GraduateStatus) (string, []float64, []float64) {
	var buffer bytes.Buffer
	var totalSubjectsScores []float64
	var totalNonSubjectsScores []float64

	buffer.WriteString("-- tb_entrance_test_factors_detail_insert" + "\n\n")

	for i := 1; i <= rows; i++ {
		var artsPhysicalSubjectsScore, generalSubjectsScore *float64
		var score1_2, score2_1, score2_2, score3_1, score3_2 *float64
		var totalNonSubjectsScore, totalSubjectsScore, volunteerScore, attendanceScore float64

		switch graduateStatuses[i-1] {
		case CANDIDATE:

			score1_2 = generateFloatPointer(0, 18)
			score2_1 = generateFloatPointer(0, 45)
			score2_2 = generateFloatPointer(0, 45)
			score3_1 = generateFloatPointer(0, 72)
			score3_2 = nil

			generalSubjectsScoreValue := *score2_1 + *score2_2 + *score3_1
			generalSubjectsScore = &generalSubjectsScoreValue
			artsPhysicalSubjectsScoreValue := *generateFloatPointer(0, 60)
			artsPhysicalSubjectsScore = &artsPhysicalSubjectsScoreValue
			totalSubjectsScore = *generalSubjectsScore + *artsPhysicalSubjectsScore

			attendanceScore = *generateFloatPointer(0, 30)
			volunteerScore = *generateFloatPointer(0, 30)
			totalNonSubjectsScore = attendanceScore + volunteerScore

		case GRADUATE:

			score1_2 = nil
			score2_1 = generateFloatPointer(0, 36)
			score2_2 = generateFloatPointer(0, 36)
			score3_1 = generateFloatPointer(0, 54)
			score3_2 = generateFloatPointer(0, 54)

			artsPhysicalSubjectsScoreValue := *generateFloatPointer(0, 60)
			artsPhysicalSubjectsScore = &artsPhysicalSubjectsScoreValue
			generalSubjectsScoreValue := *score2_1 + *score2_2 + *score3_1 + *score3_2
			generalSubjectsScore = &generalSubjectsScoreValue
			totalSubjectsScore = *generalSubjectsScore + *artsPhysicalSubjectsScore

			attendanceScore = *generateFloatPointer(0, 30)
			volunteerScore = *generateFloatPointer(0, 30)
			totalNonSubjectsScore = attendanceScore + volunteerScore

		case GED:
			artsPhysicalSubjectsScore = nil
			generalSubjectsScore = nil
			score1_2, score2_1, score2_2, score3_1, score3_2 = nil, nil, nil, nil, nil

			attendanceScore = 30                          // 검정고시는 출결점수 30점 고정
			volunteerScore = *generateFloatPointer(0, 30) //TODO: 검정고시 점수에 비례하여 봉사점수 할당
			totalNonSubjectsScore = attendanceScore + volunteerScore
			totalSubjectsScore = *generateFloatPointer(0, 240) //TODO: 검정고시 점수에 비례하여 봉사점수 할당
		}

		// 생성한 교과 성적, 비교과 성적을 배열에 저장해서 반환, tb_entrance_test_result DML을 생성할때 서류전형 총점 계산시에 사용
		totalSubjectsScores = append(totalSubjectsScores, totalSubjectsScore)
		totalNonSubjectsScores = append(totalNonSubjectsScores, totalNonSubjectsScore)

		query := fmt.Sprintf(
			"INSERT INTO tb_entrance_test_factors_detail (arts_physical_subjects_score, attendance_score, general_subjects_score, score_1_2, score_2_1, score_2_2, score_3_1, score_3_2, total_non_subjects_score, total_subjects_score, volunteer_score) "+
				"VALUES (%s, %.3f, %s, %s, %s, %s, %s, %s, %.3f, %.3f, %.3f);",
			formatNullableFloat(artsPhysicalSubjectsScore, 3),
			attendanceScore,
			formatNullableFloat(generalSubjectsScore, 3),
			formatNullableFloat(score1_2, 3),
			formatNullableFloat(score2_1, 3),
			formatNullableFloat(score2_2, 3),
			formatNullableFloat(score3_1, 3),
			formatNullableFloat(score3_2, 3),
			totalNonSubjectsScore,
			totalSubjectsScore,
			volunteerScore,
		)

		buffer.WriteString(query + "\n")
	}

	return buffer.String(), totalSubjectsScores, totalNonSubjectsScores
}
