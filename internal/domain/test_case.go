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
	TestCaseCreateInput struct {
		TaskID string `json:"task_id" db:"task_id"`
		Input  string `json:"input" db:"input"`
		Output string `json:"output" db:"output"`
	}

	TestCaseUpdateInput struct {
		Input  *string `json:"input" db:"input"`
		Output *string `json:"output" db:"output"`
	}
)

type (
	TestCaseCreateDTO struct {
		Input TestCaseCreateInput
	}

	TestCaseUpdateDTO struct {
		CaseID string
		Input  TestCaseUpdateInput
	}

	TestCaseDeleteDTO struct {
		CaseID string
	}
)
