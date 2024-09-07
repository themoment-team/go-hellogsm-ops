package main

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

var names = []string{"신희성", "김겸비", "방가온", "김재균", "김주은", "최장우", "하제우", "이승제", "전예빈", "김형록", "유시온", "양시준", "전지환", "진예원", "서채운"}
var domains = []string{"gsm.hs.kr", "gmail.com"}
var authTypes = []string{"GOOGLE", "KAKAO"}
var sexes = []string{"MALE", "FEMALE"}
var role = "APPLICANT"

func GetRandomName() string {
	return names[rand.Intn(len(names))]
}

func GetRandomEmail() string {
	domain := domains[rand.Intn(len(domains))]
	uuidPart := uuid.New().String()[:8]

	return fmt.Sprintf("%s@%s", uuidPart, domain)
}

func GetRandomPhoneNumber() string {
	return fmt.Sprintf("010%08d", rand.Intn(100000000))
}

func GetRandomDate() string {
	year := rand.Intn(10) + 2000
	month := rand.Intn(12) + 1
	day := rand.Intn(28) + 1
	return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
}

func GetCurrentTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05.000000")
}

func GenerateMemberInsertQuery(rows int) string {
	var buffer bytes.Buffer

	buffer.WriteString("-- tb_member" + "\n\n")

	for i := 1; i <= rows; i++ {
		name := GetRandomName()
		email := GetRandomEmail()
		phone := GetRandomPhoneNumber()
		birth := GetRandomDate()
		timestamp := GetCurrentTimestamp()
		authType := authTypes[rand.Intn(len(authTypes))]
		sex := sexes[rand.Intn(len(sexes))]

		query := fmt.Sprintf(
			"INSERT INTO tb_member (member_id, birth, created_time, updated_time, auth_referrer_type, email, name, phone_number, role, sex) "+
				"VALUES (%d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s');",
			i, birth, timestamp, timestamp, authType, email, name, phone, role, sex,
		)

		buffer.WriteString(query + "\n")
	}

	return buffer.String()
}
