package domain

import "lcode/pkg/db"

type (
	Comment struct {
		ID        string  `json:"id" db:"id"`
		ParentID  *string `json:"parent_id" db:"parent_id"` // parent comment id
		EntityID  string  `json:"entity_id" db:"entity_id"` // article or problem id
		Text      string  `json:"comment_text" db:"comment_text"`
		CreatedAt IntTime `json:"created_at" db:"created_at"`
		Author    `json:"author"`
	}

	Thread struct {
		Comment Comment   `json:"comment"`
		Replies []Comment `json:"replies"`
	}

	ThreadList struct {
		Threads    []Thread     `json:"threads"`
		Pagination IdPagination `json:"pagination"`
	}
)

type (
	CommentParamsInput struct {
		Sort       CommentSort        `json:"sort"`
		Pagination IdPaginationParams `json:"pagination"`
	}

	CommentSort struct {
		ByDate db.SortType `json:"date"`
	}
)

type (
	CommentCreateInput struct {
		AuthorID string  `json:"-"`
		ParentID *string `json:"parent_id"`
		EntityID string  `json:"entity_id"`
		Text     string  `json:"comment_text"`
	}

	CommentUpdateInput struct {
		ID   string  `json:"-"`
		Text *string `json:"comment_text"`
	}
)

type (
	CommentCreateDTO struct {
		Input CommentCreateInput
	}

	CommentUpdateDTO struct {
		User  User
		Input CommentUpdateInput
	}

	CommentDeleteDTO struct {
		User User
		ID   string
	}

	CommentParamsDTO struct {
		EntityID string
		Input    CommentParamsInput
	}
)
