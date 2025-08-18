package semester

import gonanoid "github.com/matoous/go-nanoid/v2"

type Semester struct {
	Semester  string `db:"semester"`
	Code      string `db:"code"`
	IsCurrent bool   `db:"is_current"`
}

func NewSemester(semester string) *Semester {

	return &Semester{
		Semester:  semester,
		Code:      generateEnrolmentCode(),
		IsCurrent: true,
	}
}

func generateEnrolmentCode() string {
	const alphabets = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const codelen = 6

	return gonanoid.MustGenerate(alphabets, codelen)
}
