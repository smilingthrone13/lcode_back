package domain

type (
	TaskTemplate struct {
		ID         string       `json:"id" db:"id"`
		TaskID     string       `json:"task_id" db:"task_id"`
		LanguageID LanguageType `json:"language_id" db:"language_id"`
		Template   string       `json:"template" db:"template"`
		Wrapper    string       `json:"wrapper" db:"wrapper"`
	}
)

type (
	TaskTemplateCreateInput struct {
		LanguageID LanguageType `json:"language_id" db:"language_id"`
		Template   string       `json:"template" db:"template"`
		Wrapper    string       `json:"wrapper" db:"wrapper"`
	}

	TaskTemplateUpdateInput struct {
		Template *string `json:"template" db:"template"`
		Wrapper  *string `json:"wrapper" db:"wrapper"`
	}
)

type (
	TaskTemplateCreateDTO struct {
		TaskID string
		Input  TaskTemplateCreateInput
	}

	TaskTemplateUpdateDTO struct {
		TaskID     string
		TemplateID string
		Input      TaskTemplateUpdateInput
	}

	TaskTemplateDeleteDTO struct {
		TemplateID string
	}
)
