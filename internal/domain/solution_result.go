package domain

type SolutionResult struct {
	SolutionID      string      `json:"solution_id" db:"solution_id"`
	TestCaseID      string      `json:"test_case_id" db:"test_case_id"`
	SubmissionToken string      `json:"-" db:"submission_token"`
	Status          JudgeStatus `json:"status" db:"status"`
	Runtime         float64     `json:"runtime" db:"runtime"`
	Memory          int         `json:"memory" db:"memory"`
	Stdout          *string     `json:"stdout" db:"stdout"`
	Stderr          *string     `json:"stderr" db:"stderr"`
}

type GetSolutionResultsDTO struct {
	SolutionID string
	User       User
}

func (d GetSolutionResultsDTO) GetSolutionID() string {
	return d.SolutionID
}

func (d GetSolutionResultsDTO) GetUser() User {
	return d.User
}
