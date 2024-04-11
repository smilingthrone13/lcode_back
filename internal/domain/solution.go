package domain

type SolutionStatus string

const (
	SolutionStatusTesting   SolutionStatus = "testing"
	SolutionStatusCompleted SolutionStatus = "completed"
	SolutionStatusError     SolutionStatus = "error"
)

type Solution struct {
	Id         string         `json:"id" db:"id"`
	UserID     string         `json:"user_id" db:"user_id"`
	TaskID     string         `json:"task_id" db:"task_id"`
	LanguageID LanguageType   `json:"language_id" db:"language_id"`
	Code       string         `json:"-" db:"code"`
	Status     SolutionStatus `json:"status" db:"status"`
	Runtime    float64        `json:"runtime" db:"runtime"`
	Memory     int            `json:"memory" db:"memory"`
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
type IGetSolutionDTO interface {
	GetSolutionID() string
	GetUser() User
}

type CreateSolutionDTO struct {
	TaskID     string
	LanguageID LanguageType
	Code       string
	User       User
}

type GetSolutionsDTO struct {
	TaskID string
	User   User
}

type GetSolutionCodeDTO struct {
	SolutionID string
	User       User
}

func (d GetSolutionCodeDTO) GetSolutionID() string {
	return d.SolutionID
}

func (d GetSolutionCodeDTO) GetUser() User {
	return d.User
}

type UpdateSolutionDTO struct {
	ID      string
	Status  *SolutionStatus
	Runtime *float64
	Memory  *int
}
