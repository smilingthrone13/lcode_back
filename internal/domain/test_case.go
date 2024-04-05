package domain

type (
	TestCase struct {
		ID     string                 `json:"id" db:"id"`
		TaskID string                 `json:"task_id" db:"task_id"`
		Number string                 `json:"number" db:"number"`
		Input  map[string]interface{} `json:"input" db:"input"`
		Output map[string]interface{} `json:"output" db:"output"`
	}
)
