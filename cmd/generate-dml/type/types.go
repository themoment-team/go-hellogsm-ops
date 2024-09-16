package _type

type GraduateStatus string
type Screening string
type OneseoStatus string

type ScreeningParam string

const (
	CANDIDATE              GraduateStatus = "CANDIDATE"
	GRADUATE               GraduateStatus = "GRADUATE"
	GED                    GraduateStatus = "GED"
	RANDOM_GRADUATE_STATUS GraduateStatus = "RANDOM"
)

const (
	GENERAL         Screening = "GENERAL"
	SPECIAL         Screening = "SPECIAL"
	EXTRA_ADMISSION Screening = "EXTRA_ADMISSION"
	EXTRA_VETERANS  Screening = "EXTRA_VETERANS"
)

const (
	FIRST       OneseoStatus = "FIRST"
	SECOND      OneseoStatus = "SECOND"
	FINAL_MAJOR OneseoStatus = "FINAL_MAJOR"
	RE_EVALUATE OneseoStatus = "RE_EVALUATE"
)

const (
	GEN ScreeningParam = "GEN"
	SPE ScreeningParam = "SPE"
	EXT ScreeningParam = "EXT"
)

func (g GraduateStatus) IsValidGraduateStatus() bool {
	switch g {
	case CANDIDATE, GRADUATE, GED, RANDOM_GRADUATE_STATUS:
		return true
	}
	return false
}

func (t OneseoStatus) IsValidOneseoStatus() bool {
	switch t {
	case FIRST, SECOND, FINAL_MAJOR, RE_EVALUATE:
		return true
	}
	return false
}
