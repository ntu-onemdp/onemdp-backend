package semester

import (
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

type SemesterService struct {
	semesterRepo *SemesterRepository

	// Cache the current semester and code for easy retrieval.
	currentSemester *string
	code            *string
}

var Service *SemesterService

func NewSemesterService(semRepo *SemesterRepository) *SemesterService {
	semester, err := semRepo.getCurrentSem()
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error initializing semester service")
		return nil
	}

	return &SemesterService{
		semesterRepo:    semRepo,
		code:            &semester.Code,
		currentSemester: &semester.Semester,
	}
}

// Get the current sem
func (s *SemesterService) GetCurrentSem() *string {
	return s.currentSemester
}

// Get the enrolment code for the current sem
func (s *SemesterService) GetCode() string {
	return *s.code
}

// Insert a new semester into the system
func (s *SemesterService) NewSemester(sem string) error {
	semester := NewSemester(sem)

	// Insert new semester into the repo
	if err := s.semesterRepo.insert(*semester); err != nil {
		utils.Logger.Error().Err(err).Msg("Error inserting new semester into repo")
		return err
	}

	// Update cache
	s.code = &semester.Code
	s.currentSemester = &semester.Semester
	utils.Logger.Trace().Str("current semester", *s.currentSemester).Str("enrolment code", *s.code).Msg("Semester service cache updated successfully")

	return nil
}

// Refresh code for current semester
// Returns new code if successful
func (s *SemesterService) RefreshCode() (string, error) {
	// Generate new code
	code := generateEnrolmentCode()

	// Update repo with new enrolment code
	if err := s.semesterRepo.RefreshCode(code); err != nil {
		utils.Logger.Error().Err(err).Msg("Error updating database with new enrolment code")
		return "", err
	}

	// Update cache
	s.code = &code
	utils.Logger.Trace().Str("enrolment code", code).Msg("Enrolment code updated in cache")
	utils.Logger.Info().Str("enrolment code", code).Msg("Enrolment code updated successfully.")

	return code, nil
}
