package domain

type SolutionStatus string

const (
	Testing   SolutionStatus = "testing"
	Completed SolutionStatus = "completed"
	Error     SolutionStatus = "error"
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

type CreateSolutionDTO struct {
	UserID     string
	TaskID     string
	LanguageID LanguageType
	Code       string
	Status     SolutionStatus
	Runtime    float64
	Memory     int
}

type UpdateSolutionDTO struct {
	ID      string
	Status  *SolutionStatus
	Runtime *float64
	Memory  *int
}
