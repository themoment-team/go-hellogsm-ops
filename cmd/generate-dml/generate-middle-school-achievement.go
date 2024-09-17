package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

func randomScoreArray(length int, min int, max int) string {
	arr := make([]string, length)
	for i := 0; i < length; i++ {
		arr[i] = strconv.Itoa(rand.Intn(max-min+1) + min)
	}
	return "[" + strings.Join(arr, ", ") + "]"
}

func GenerateMiddleSchoolAchievementInsertQuery(rows int, graduateStatuses []GraduateStatus) string {
	var buffer bytes.Buffer

	buffer.WriteString("-- tb_middle_school_achievement_insert" + "\n\n")

	for i := 1; i <= rows; i++ {

		artsPhysicalSubjects := `'["체육","미술","음악"]'`
		generalSubjects := `'["국어","도덕","사회","역사","수학","과학","기술가정","영어"]'`
		newSubjects := `'["프로그래밍"]'`

		graduateStatus := graduateStatuses[i-1]
		var gedTotalScore string
		var absentDays, achievement_2_1, achievement_2_2, achievement_3_1, achievement_3_2, artsPhysicalAchievement, attendanceDays, volunteerTime string
		var freeSemester, liberalSystem string

		switch graduateStatus {
		case CANDIDATE:
			gedTotalScore = "NULL"
			absentDays = fmt.Sprintf("%s%v%s", "'", randomScoreArray(3, 0, 3), "'")
			achievement_2_1 = fmt.Sprintf("%s%v%s", "'", randomScoreArray(9, 1, 5), "'")
			achievement_2_2 = fmt.Sprintf("%s%v%s", "'", randomScoreArray(9, 1, 5), "'")
			achievement_3_1 = fmt.Sprintf("%s%v%s", "'", randomScoreArray(9, 1, 5), "'")
			achievement_3_2 = "NULL"
			artsPhysicalAchievement = fmt.Sprintf("%s%v%s", "'", randomScoreArray(9, 3, 5), "'")
			attendanceDays = fmt.Sprintf("%s%v%s", "'", randomScoreArray(9, 1, 5), "'")
			volunteerTime = fmt.Sprintf("%s%v%s", "'", randomScoreArray(3, 0, 5), "'")
			freeSemester = "NULL"
			liberalSystem = "'자유학년제'"

		case GRADUATE:
			gedTotalScore = "NULL"
			absentDays = fmt.Sprintf("%s%v%s", "'", randomScoreArray(3, 0, 3), "'")
			achievement_2_1 = fmt.Sprintf("%s%v%s", "'", randomScoreArray(9, 1, 5), "'")
			achievement_2_2 = fmt.Sprintf("%s%v%s", "'", randomScoreArray(9, 1, 5), "'")
			achievement_3_1 = fmt.Sprintf("%s%v%s", "'", randomScoreArray(9, 1, 5), "'")
			achievement_3_2 = fmt.Sprintf("%s%v%s", "'", randomScoreArray(9, 1, 5), "'")
			artsPhysicalAchievement = fmt.Sprintf("%s%v%s", "'", randomScoreArray(12, 3, 5), "'")
			attendanceDays = fmt.Sprintf("%s%v%s", "'", randomScoreArray(9, 1, 5), "'")
			volunteerTime = fmt.Sprintf("%s%v%s", "'", randomScoreArray(3, 0, 5), "'")
			freeSemester = "NULL"
			liberalSystem = "'자유학년제'"

		case GED:
			gedTotalScore = fmt.Sprintf("%s%d%s", "'", rand.Intn(201)+400, "'")
			absentDays = "NULL"
			achievement_2_1 = "NULL"
			achievement_2_2 = "NULL"
			achievement_3_1 = "NULL"
			achievement_3_2 = "NULL"
			attendanceDays = "NULL"
			artsPhysicalAchievement = "NULL"
			volunteerTime = "NULL"
			freeSemester = "NULL"
			liberalSystem = "NULL"
			artsPhysicalSubjects = "NULL"
			generalSubjects = "NULL"
			newSubjects = "NULL"
		}

		query := fmt.Sprintf(
			"INSERT INTO tb_middle_school_achievement (oneseo_id, ged_total_score, absent_days, achievement_1_2, achievement_2_1, achievement_2_2, achievement_3_1, achievement_3_2, arts_physical_achievement, arts_physical_subjects, attendance_days, free_semester, general_subjects, liberal_system, new_subjects, volunteer_time) "+
				"VALUES (%d, %s, %s, NULL, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s);",
			i, gedTotalScore, absentDays, achievement_2_1, achievement_2_2, achievement_3_1, achievement_3_2, artsPhysicalAchievement, artsPhysicalSubjects,
			attendanceDays, freeSemester, generalSubjects, liberalSystem, newSubjects, volunteerTime,
		)

		buffer.WriteString(query + "\n")
	}

	return buffer.String()
}
