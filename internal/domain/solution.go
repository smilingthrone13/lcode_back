package domain

type SolutionStatus string

const (
	SolutionStatusTesting   SolutionStatus = "testing"
	SolutionStatusCompleted SolutionStatus = "completed"
	SolutionStatusError     SolutionStatus = "error"
)

type Solution struct {
	Id         string         `json:"id"`
	UserID     string         `json:"user_id"`
	TaskID     string         `json:"task_id"`
	LanguageID LanguageType   `json:"language_id"`
	Code       string         `json:"code"`
	Status     SolutionStatus `json:"status"`
	Runtime    float64        `json:"runtime"`
	Memory     int            `json:"memory"`
}

// entity
type CreateSolutionEntity struct {
	TaskID     string
	LanguageID LanguageType
	Code       string
	Status     SolutionStatus
	User       User
}

// dto
type CreateSolutionDTO struct {
	TaskID     string
	LanguageID LanguageType
	Code       string
	User       User
}

type UpdateSolutionDTO struct {
	ID      string
	Status  *SolutionStatus
	Runtime *float64
	Memory  *int
}
