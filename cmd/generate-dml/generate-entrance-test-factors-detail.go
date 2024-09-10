package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"themoment-team/hellogsm-notice-server/cmd/generate-dml/type"
)

func randomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func GenerateEntranceTestFactorsDetailInsertQuery(rows int, graduateStatuses []_type.GraduateStatus) (string, []float64, []float64) {
	var buffer bytes.Buffer
	var totalSubjectsScores []float64
	var totalNonSubjectsScores []float64

	buffer.WriteString("-- tb_entrance_test_factors_detail" + "\n\n")

	for i := 1; i <= rows; i++ {
		var artsPhysicalSubjectsScore, attendanceScore, generalSubjectsScore *float64
		var score1_2, score2_1, score2_2, score3_1, score3_2 *float64
		var totalNonSubjectsScore, totalSubjectsScore, volunteerScore float64

		switch graduateStatuses[i-1] {
		case _type.CANDIDATE:
			score2_1Value := randomFloat(0, 54)
			score2_2Value := randomFloat(0, 54)
			score3_1Value := randomFloat(0, 72)

			score1_2 = nil
			score2_1 = &score2_1Value
			score2_2 = &score2_2Value
			score3_1 = &score3_1Value
			score3_2 = nil

			generalSubjectsScoreValue := *score2_1 + *score2_2 + *score3_1
			generalSubjectsScore = &generalSubjectsScoreValue
			artsPhysicalSubjectsScoreValue := randomFloat(0, 60)
			artsPhysicalSubjectsScore = &artsPhysicalSubjectsScoreValue
			totalSubjectsScore = *generalSubjectsScore + *artsPhysicalSubjectsScore

			attendanceScore = randomFloatPointer(0, 30)
			volunteerScore = randomFloat(0, 30)
			totalNonSubjectsScore = *attendanceScore + volunteerScore

		case _type.GRADUATE:
			score2_1Value := randomFloat(0, 32)
			score2_2Value := randomFloat(0, 32)
			score3_1Value := randomFloat(0, 54)
			score3_2Value := randomFloat(0, 54)

			score1_2 = nil
			score2_1 = &score2_1Value
			score2_2 = &score2_2Value
			score3_1 = &score3_1Value
			score3_2 = &score3_2Value

			artsPhysicalSubjectsScoreValue := randomFloat(0, 60)
			artsPhysicalSubjectsScore = &artsPhysicalSubjectsScoreValue
			generalSubjectsScoreValue := *score2_1 + *score2_2 + *score3_1 + *score3_2
			generalSubjectsScore = &generalSubjectsScoreValue
			totalSubjectsScore = *generalSubjectsScore + *artsPhysicalSubjectsScore

			attendanceScore = randomFloatPointer(0, 30)
			volunteerScore = randomFloat(0, 30)
			totalNonSubjectsScore = *attendanceScore + volunteerScore

		case _type.GED:
			artsPhysicalSubjectsScore = nil
			generalSubjectsScore = nil
			score1_2, score2_1, score2_2, score3_1, score3_2 = nil, nil, nil, nil, nil

			attendanceScore = randomFloatPointer(0, 30)
			volunteerScore = randomFloat(0, 30)
			totalNonSubjectsScore = *attendanceScore + volunteerScore
			totalSubjectsScore = randomFloat(0, 240)
		}

		totalSubjectsScores = append(totalSubjectsScores, totalSubjectsScore)
		totalNonSubjectsScores = append(totalNonSubjectsScores, totalNonSubjectsScore)

		query := fmt.Sprintf(
			"INSERT INTO tb_entrance_test_factors_detail (arts_physical_subjects_score, attendance_score, general_subjects_score, score_1_2, score_2_1, score_2_2, score_3_1, score_3_2, total_non_subjects_score, total_subjects_score, volunteer_score) "+
				"VALUES (%s, %.2f, %s, %s, %s, %s, %s, %s, %.2f, %.2f, %.2f);",
			formatNullable(artsPhysicalSubjectsScore),
			*attendanceScore,
			formatNullable(generalSubjectsScore),
			formatNullable(score1_2),
			formatNullable(score2_1),
			formatNullable(score2_2),
			formatNullable(score3_1),
			formatNullable(score3_2),
			totalNonSubjectsScore,
			totalSubjectsScore,
			volunteerScore,
		)

		buffer.WriteString(query + "\n")
	}

	return buffer.String(), totalSubjectsScores, totalNonSubjectsScores
}

func randomFloatPointer(min, max float64) *float64 {
	value := randomFloat(min, max)
	return &value
}

func formatNullable(value *float64) string {
	if value == nil {
		return "NULL"
	}
	return fmt.Sprintf("%.2f", *value)
}
