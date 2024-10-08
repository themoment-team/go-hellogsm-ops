package main

import (
	"bytes"
	"fmt"
)

func GenerateOneseoPrivacyDetailInsertQuery(rows int, graduateStatuses []GraduateStatus) string {
	var buffer bytes.Buffer

	buffer.WriteString("-- tb_oneseo_privacy_detail_insert" + "\n\n")

	for i := 1; i <= rows; i++ {
		address := "'광주광역시 광산구 상무대로 312'"
		detailAddress := "'101-1001'"
		guardianName := "'김보호'"
		guardianPhoneNumber := fmt.Sprintf("%s%s%s", "'", GetRandomPhoneNumber(), "'")
		profileImg := "'https://asd'"
		relationshipWithGuardian := "'모'"

		schoolAddress := "'광주광역시 광산구 상무대로 312'"
		schoolName := "'광주소프트웨어마이스터중학교'"
		schoolTeacherName := "'김선생'"
		schoolTeacherPhoneNumber := fmt.Sprintf("%s%s%s", "'", GetRandomPhoneNumber(), "'")
		graduationDate := fmt.Sprintf("%s%s%s", "'", GetRandomDate()[:7], "'")

		// GraduateStatus가 GED라면 중학교 관련 정보에 NULL값 할당
		graduateStatus := graduateStatuses[i-1]
		if graduateStatus == GED {
			schoolAddress = "NULL"
			schoolName = "NULL"
			schoolTeacherName = "NULL"
			schoolTeacherPhoneNumber = "NULL"
		}

		query := fmt.Sprintf(
			"INSERT INTO tb_oneseo_privacy_detail (oneseo_id, address, detail_address, graduation_type, guardian_name, guardian_phone_number, profile_img, relationship_with_guardian, school_address, school_name, school_teacher_name, school_teacher_phone_number, graduation_date) "+
				"VALUES (%d, %s, %s, '%s', %s, %s, %s, %s, %s, %s, %s, %s, %s);",
			i, address, detailAddress, graduateStatus, guardianName, guardianPhoneNumber, profileImg, relationshipWithGuardian,
			schoolAddress, schoolName, schoolTeacherName, schoolTeacherPhoneNumber, graduationDate,
		)

		buffer.WriteString(query + "\n")
	}

	return buffer.String()
}
