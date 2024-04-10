package domain

import "lcode/pkg/db"

type (
	Task struct {
		ID           string  `json:"id" db:"id"`
		Number       string  `json:"number" db:"number"`
		Name         string  `json:"name" db:"name"`
		Description  string  `json:"description" db:"description"`
		Category     string  `json:"category" db:"category"`
		Difficulty   string  `json:"difficulty" db:"difficulty"`
		RuntimeLimit float64 `json:"runtime_limit" db:"runtime_limit"`
		MemoryLimit  int     `json:"memory_limit" db:"memory_limit"`
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

	TaskAttributes struct {
		Categories   []string `json:"categories" db:"categories"`
		Difficulties []string `json:"difficulties" db:"difficulties"`
	}
)

type (
	TaskCreateInput struct {
		Name         string  `json:"name" db:"name"`
		Description  *string `json:"description" db:"description"`
		Category     string  `json:"category" db:"category"`
		Difficulty   string  `json:"difficulty" db:"difficulty"`
		RuntimeLimit *string `json:"runtime_limit" db:"runtime_limit"`
		MemoryLimit  *string `json:"memory_limit" db:"memory_limit"`
	}

	TaskUpdateInput struct {
		Name         *string `json:"name" db:"name"`
		Description  *string `json:"description" db:"description"`
		Category     *string `json:"category" db:"category"`
		Difficulty   *string `json:"difficulty" db:"difficulty"`
		RuntimeLimit *string `json:"runtime_limit" db:"runtime_limit"`
		MemoryLimit  *string `json:"memory_limit" db:"memory_limit"`
	}

	TaskParamsInput struct {
		Sort       TaskSort           `json:"sort"`
		Pagination IdPaginationParams `json:"pagination"`
	}
)

type (
	TaskUpdateDTO struct {
		TaskID string
		Input  TaskUpdateInput
	}

	TaskParamsDTO struct {
		Input TaskParams
	}
)
