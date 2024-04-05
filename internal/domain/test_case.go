package domain

type (
	TestCase struct {
		ID     string `json:"id" db:"id"`
		Number string `json:"number" db:"number"`
		TaskID string `json:"task_id" db:"task_id"`
		Input  string `json:"input" db:"input"`
		Output string `json:"output" db:"output"`
	}
)

type (
	TestCaseCreate struct {
		TaskID string `json:"task_id" db:"task_id"`
		Input  string `json:"input" db:"input"`
		Output string `json:"output" db:"output"`
	}

	TestCaseUpdate struct {
		Input  *string `json:"input" db:"input"`
		Output *string `json:"output" db:"output"`
	}
)
