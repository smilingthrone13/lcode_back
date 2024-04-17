package domain

import "lcode/pkg/db"

type CommentOriginType string

const (
	ArticleOriginType CommentOriginType = "article_comment"
	TaskOriginType    CommentOriginType = "task_comment"
)

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
		Sort       CommentSort
		Pagination IdPaginationParams
	}

	CommentSort struct {
		ByDate db.SortType
	}
)

type (
	CommentCreateInput struct {
		AuthorID string  `json:"-"`
		EntityID string  `json:"-"`
		ParentID *string `json:"parent_id"`
		Text     string  `json:"comment_text"`
	}

	CommentUpdateInput struct {
		ID   string  `json:"-"`
		Text *string `json:"comment_text"`
	}
)

type (
	CommentCreateDTO struct {
		OriginType CommentOriginType
		Input      CommentCreateInput
	}

	CommentUpdateDTO struct {
		User       User
		OriginType CommentOriginType
		Input      CommentUpdateInput
	}

	CommentDeleteDTO struct {
		User       User
		OriginType CommentOriginType
		ID         string
	}

	CommentParamsDTO struct {
		EntityID   string
		OriginType CommentOriginType
		Input      CommentParamsInput
	}
)
