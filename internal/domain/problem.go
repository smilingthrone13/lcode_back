package domain

type (
	Problem struct {
		Task          Task           `json:"task"`
		TaskTemplates []TaskTemplate `json:"task_templates"`
		TestCases     []TestCase     `json:"test_cases"`
	}
)

type (
	ProblemCreateInput struct {
		Task          TaskCreateInput           `json:"task"`
		TaskTemplates []TaskTemplateCreateInput `json:"task_templates"`
		TestCases     []TestCaseCreateInput     `json:"test_cases"`
	}
)

type (
	ProblemCreateDTO struct {
		Input ProblemCreateInput
	}

	ProblemDeleteDTO struct {
		TaskID string
	}

	GetProblemDTO struct {
		TaskID string
	}
)
