package domain

type SolutionResult struct {
	SolutionID      string      `json:"solution_id"`
	TestCaseID      string      `json:"test_case_id"`
	SubmissionToken string      `json:"submission_token"`
	Status          JudgeStatus `json:"status"`
	Runtime         float64     `json:"runtime"`
	Memory          int         `json:"memory"`
}
