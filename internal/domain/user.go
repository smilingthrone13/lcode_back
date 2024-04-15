package domain

import (
	"io"
)

const (
	UserCtxKey = "user"
	DtoCtxKey  = "dto"

	UsersFolder     = "users"
	AvatarFolder    = "avatar"
	ThumbnailFolder = "thumbnail"
)

type (
	User struct {
		ID           string `json:"id" db:"id" mapstructure:"id"`
		Email        string `json:"email" db:"email" mapstructure:"email"`
		Username     string `json:"username" db:"username" mapstructure:"username"`
		FirstName    string `json:"first_name" db:"first_name" mapstructure:"first_name"`
		LastName     string `json:"last_name" db:"last_name" mapstructure:"last_name"`
		IsAdmin      bool   `json:"is_admin" db:"is_admin" mapstructure:"is_admin"`
		PasswordHash string `json:"-" db:"password_hash"`
	}

	Author struct {
		UserID    string `json:"user_id" db:"user_id"`
		Username  string `json:"username" db:"username"`
		FirstName string `json:"first_name" db:"first_name"`
		LastName  string `json:"last_name" db:"last_name"`
	}
)

// entity for input for repository
type (
	CreateUserEntity struct {
		Email        string
		Username     string
		FirstName    string
		LastName     string
		PasswordHash string
	}

	UpdateUserEntity struct {
		UserID       string
		Email        *string
		Username     *string
		FirstName    *string
		LastName     *string
		PasswordHash *string
		IsAdmin      *bool
	}
)

// layer transfer objects
type (
	CreateUserDTO struct {
		Email     string `json:"email" binding:"required,email,max=100"`
		Username  string `json:"username" binding:"required,min=3,max=50"`
		FirstName string `json:"first_name" binding:"required,min=2,max=50"`
		LastName  string `json:"last_name" binding:"required,min=2,max=50"`
		Password  string `json:"password" binding:"required,min=5,max=50"`
	}

	LoginDTO struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	UploadUserAvatarDTO struct {
		Media        io.Reader
		FullFileName string
		Name         string
		Extension    string
		MediaType    string
		User         User
	}

	DeleteUserAvatarDTO struct {
		User User
	}

	RefreshTokenDTO struct {
		RefreshToken string `json:"refresh_token"`
	}

	UpdateUserDTO struct {
		UserID    string
		Email     *string `json:"email" binding:"email,max=100"`
		Username  *string `json:"username" binding:"min=3,max=50"`
		FirstName *string `json:"first_name" binding:"min=2,max=50"`
		LastName  *string `json:"last_name" binding:"min=2,max=50"`
		Password  *string `json:"password" binding:"min=5,max=50"`
		IsAdmin   *bool   `json:"is_admin"`
	}

	ChangeUserAdminPermissionDTO struct {
		UserID  string `json:"user_id"`
		IsAdmin bool   `json:"is_admin"`
	}

	ChangeUserPasswordDTO struct {
		UserID      string `json:"user_id"`
		NewPassword string `json:"new_password"`
	}
)
