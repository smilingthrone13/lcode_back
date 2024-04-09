package domain

import "lcode/pkg/db"

type (
	Comment struct {
		ID        string  `json:"id"`
		AuthorID  string  `json:"user_id"`
		ParentID  string  `json:"parent_id"` // parent comment id
		EntityID  string  `json:"post_id"`   // article or problem id
		Content   string  `json:"content"`
		CreatedAt IntTime `json:"created_at"`
	}

	Thread struct {
		EntityID string    `json:"entity_id"`
		Comment  Comment   `json:"comment"`
		Replies  []Comment `json:"replies"`
	}

	ThreadList struct {
		Threads    []Thread     `json:"threads"`
		Pagination IdPagination `json:"pagination"`
	}
)

type (
	CommentParams struct {
		Sort       CommentSort
		Pagination IdPaginationParams
	}

	CommentSort struct {
		ByDate db.SortType `json:"date"`
	}
)

type (
	CommentCreateInput struct {
		AuthorID string  `json:"-"`
		ParentID *string `json:"parent_id"`
		EntityID string  `json:"post_id"`
		Content  string  `json:"content"`
	}

	CommentUpdateInput struct {
		ID      string  `json:"-"`
		Content *string `json:"content"`
	}

	CommentDeleteInput struct {
		ID string
	}
)

type (
	CommentCreateDTO struct {
		Input CommentCreateInput
	}

	CommentUpdateDTO struct {
		Input CommentUpdateInput
	}

	CommentDeleteDTO struct {
		Input CommentDeleteInput
	}

	CommentParamsDTO struct {
		Input CommentParams
	}
)
