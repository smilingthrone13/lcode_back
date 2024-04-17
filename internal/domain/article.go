package domain

import "lcode/pkg/db"

const PracticeArticleID = "00000000-0000-0000-0000-000000000000"

type (
	Article struct {
		ID         string   `json:"id" db:"id"`
		Title      string   `json:"title" db:"title"`
		Content    string   `json:"content" db:"content"`
		Categories []string `json:"categories" db:"categories"`
		CreatedAt  IntTime  `json:"created_at" db:"created_at"`
		Author     `json:"author"`
	}

	ArticleList struct {
		Articles   []Article    `json:"articles"`
		Pagination IdPagination `json:"pagination"`
	}
)

type (
	ArticleParams struct {
		Filter     ArticleFilter
		Sort       ArticleSort
		Pagination IdPaginationParams
	}

	ArticleFilter struct {
		Search     string
		Categories []string
	}

	ArticleSort struct {
		ByDate db.SortType
	}

	ArticleAttributes struct {
		Categories []string `json:"categories"`
	}
)

type (
	ArticleCreateInput struct {
		AuthorID   string   `json:"-"`
		Title      string   `json:"title"`
		Content    string   `json:"content"`
		Categories []string `json:"categories"`
	}

	ArticleUpdateInput struct {
		ID         string   `json:"-"`
		Title      *string  `json:"title"`
		Content    *string  `json:"content"`
		Categories []string `json:"categories"`
	}

	ArticleParamsInput struct {
		Sort       ArticleSort
		Pagination IdPaginationParams
	}
)

type (
	ArticleCreateDTO struct {
		Input ArticleCreateInput
	}

	ArticleUpdateDTO struct {
		Input ArticleUpdateInput
	}

	ArticleDeleteDTO struct {
		ID string `json:"id"`
	}

	ArticleParamsDTO struct {
		Input ArticleParams
	}

	ArticleGetDTO struct {
		ID string `json:"id"`
	}
)
