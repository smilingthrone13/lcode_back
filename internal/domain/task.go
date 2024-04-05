package domain

import "lcode/pkg/db"

type (
	Task struct {
		ID           string `json:"id" db:"id"`
		Number       string `json:"number" db:"number"`
		Name         string `json:"name" db:"name"`
		Description  string `json:"description" db:"description"`
		Category     string `json:"category" db:"category"`
		Difficulty   string `json:"difficulty" db:"difficulty"`
		RuntimeLimit string `json:"runtime_limit" db:"runtime_limit"`
		MemoryLimit  string `json:"memory_limit" db:"memory_limit"`
	}

	TaskList struct {
		Tasks      []Task       `json:"tasks"`
		Pagination IdPagination `json:"pagination"`
	}
)

type (
	TaskParams struct {
		Filter     TaskFilter
		Sort       TaskSort
		Pagination IdPaginationParams
	}

	TaskFilter struct {
		Search       string
		Categories   []string
		Difficulties []string
	}

	TaskSort struct {
		ByNumber db.SortType `json:"number"`
		// todo: ByDifficulty - за деньги да
	}

	IdPaginationParams struct {
		Limit   int     `json:"limit"`
		AfterID *string `json:"after_id"`
	}
)

type (
	TaskCreate struct {
		Name         string  `json:"name" db:"name"`
		Description  *string `json:"description" db:"description"`
		Category     string  `json:"category" db:"category"`
		Difficulty   string  `json:"difficulty" db:"difficulty"`
		RuntimeLimit *string `json:"runtime_limit" db:"runtime_limit"`
		MemoryLimit  *string `json:"memory_limit" db:"memory_limit"`
	}

	TaskUpdate struct {
		Name         *string `json:"name" db:"name"`
		Description  *string `json:"description" db:"description"`
		Category     *string `json:"category" db:"category"`
		Difficulty   *string `json:"difficulty" db:"difficulty"`
		RuntimeLimit *string `json:"runtime_limit" db:"runtime_limit"`
		MemoryLimit  *string `json:"memory_limit" db:"memory_limit"`
	}
)
