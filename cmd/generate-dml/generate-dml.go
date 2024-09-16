package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"themoment-team/hellogsm-notice-server/cmd/generate-dml/type"
)

func main() {
	graduateStatusParam := flag.String("graduate", "RANDOM", "Status of grade")
	screeningParam := flag.String("screening", "필수 입력 파라미터 입니다.", "Status of screening")
	oneseoStatusParam := flag.String("status", "필수 입력 파라미터 입니다.", "Status of oneseo")

	flag.Parse()

	graduateStatus := _type.GraduateStatus(strings.ToUpper(*graduateStatusParam))
	oneseoStatus := _type.OneseoStatus(strings.ToUpper(*oneseoStatusParam))

	generalCount := 0
	specialCount := 0
	extraCount := 0

	rows := initScreeningCount(screeningParam, &generalCount, &specialCount, &extraCount)

	err := validateParameter(graduateStatus, oneseoStatus)
	if err != nil {
		fmt.Println(err)
		return
	}

	// GraduateStatus 가 RANDOM_GRADUATE_STATUS 라면 랜덤한 GraduateStatus 배열을 생성한 후 같은 인덱스의 row 들에 공통적으로 적용
	graduateStatuses := resolveGraduateStatuses(graduateStatus, rows)

	memberInsertQuery := GenerateMemberInsertQuery(rows)
	oneseoInsertQuery := GenerateOneseoInsertQuery(rows, generalCount, specialCount, extraCount, oneseoStatus)
	oneseoPrivacyDetailInsertQuery := GenerateOneseoPrivacyDetailInsertQuery(rows, graduateStatuses)
	middleSchoolAchievementInsertQuery := GenerateMiddleSchoolAchievementInsertQuery(rows, graduateStatuses)
	generateEntranceTestFactorsDetailInsertQuery, totalSubjectsScores, totalNonSubjectsScores := GenerateEntranceTestFactorsDetailInsertQuery(rows, graduateStatuses)
	generateEntranceTestResultInsertQuery := GenerateEntranceTestResultInsertQuery(rows, oneseoStatus, totalSubjectsScores, totalNonSubjectsScores)

	createSqlFiles(
		memberInsertQuery,
		oneseoInsertQuery,
		oneseoPrivacyDetailInsertQuery,
		middleSchoolAchievementInsertQuery,
		generateEntranceTestFactorsDetailInsertQuery,
		generateEntranceTestResultInsertQuery,
	)
}

func initScreeningCount(screeningParam *string, generalCount *int, specialCount *int, extraCount *int) int {
	screnningValues := strings.Split(*screeningParam, ",")

	var rows int
	for _, s := range screnningValues {

		screeningType := s[:3]
		screeningCount := s[3:]

		count, err := strconv.Atoi(screeningCount)
		if err != nil {
			panic("전형 지원자 수를 정수로 변환하는 중 오류 발생")
		}

		switch _type.ScreeningParam(screeningType) {
		case _type.GEN:
			*generalCount += count
		case _type.SPE:
			*specialCount += count
		case _type.EXT:
			*extraCount += count
		}

		rows += count
	}

	return rows
}

func validateParameter(graduateStatus _type.GraduateStatus, oneseoStatus _type.OneseoStatus) error {
	if !graduateStatus.IsValidGraduateStatus() {
		return errors.New(fmt.Sprintf("잘못된 졸업상태가 입력되었습니다: %s", graduateStatus))
	}

	if !oneseoStatus.IsValidOneseoStatus() {
		return errors.New(fmt.Sprintf("잘못된 원서상태가 입력되었습니다: %s", oneseoStatus))
	}

	return nil
}

func resolveGraduateStatuses(graduateStatus _type.GraduateStatus, rows int) []_type.GraduateStatus {

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

	return graduateStatuses
}

func createSqlFiles(queries ...string) {
	mkdirResultIfNotExist()

	for idx, query := range queries {
		createSqlFileAndWrite(idx, query)
	}
}

func createSqlFileAndWrite(order int, query string) {
	file := createSqlFile(order, query)
	writer := bufio.NewWriter(file)
	n, err := writer.WriteString(query)
	check(err)
	fmt.Printf("wrote %d bytes\n", n)
	err = writer.Flush()
	check(err)
}

func createSqlFile(order int, query string) *os.File {
	firstLine := resolveFileName(order, query)
	file, err := os.Create(fmt.Sprintf("./result/%s.sql", firstLine))
	check(err)
	return file
}

// order 는 파일의 반영, 정렬 순서를 위해 받는다.
// query 는 파일에 쓸 sql 쿼리를 받는다.
func resolveFileName(order int, query string) string {
	lines := strings.Split(query, "\n")
	fName := lines[0]

	// 생성된 쿼리의 첫번째 라인 주석(-- tb_...)을 기반으로 fileName을 결정한다.
	m, err := regexp.MatchString("^--", fName)
	check(err)
	if m == false {
		panic("SQL 생성시, 주석(-- tb_...)이 포함되어야 함.")
	}

	// 파일명에서 주석 표기는 제거한다.
	fName = strings.ReplaceAll(fName, "-- ", "")

	// 연관관계등 순서를 보장해야하는 경우가 있어 파일에 순서를 명시한다.
	fName = fmt.Sprintf("%d_%s", order, fName)

	return fName
}

func mkdirResultIfNotExist() {
	outPath := "./result"
	fMode := 0700

	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		err := os.Mkdir(outPath, os.FileMode(fMode))
		check(err)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
