package domain

type (
	TaskTemplate struct {
		ID         string `json:"id" db:"id"`
		TaskID     string `json:"task_id" db:"task_id"`
		LanguageID int    `json:"language_id" db:"language_id"`
		Template   string `json:"template" db:"template"`
		Wrapper    string `json:"wrapper" db:"wrapper"`
	}
)

type (
	TaskTemplateCreate struct {
		TaskID     string `json:"task_id" db:"task_id"`
		LanguageID int    `json:"language_id" db:"language_id"`
		Template   string `json:"template" db:"template"`
		Wrapper    string `json:"wrapper" db:"wrapper"`
	}

	TaskTemplateUpdate struct {
		Template *string `json:"template" db:"template"`
		Wrapper  *string `json:"wrapper" db:"wrapper"`
	}
)
