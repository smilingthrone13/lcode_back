package domain

const (
	UserCtxKey = "user"
	DtoCtxKey  = "dto"
)

type User struct {
	ID           string `json:"id" mapstructure:"id"`
	Username     string `json:"username" mapstructure:"username"`
	FirstName    string `json:"first_name" mapstructure:"first_name"`
	LastName     string `json:"last_name" mapstructure:"last_name"`
	IsAdmin      bool   `json:"is_admin" mapstructure:"is_admin"`
	PasswordHash string `json:"-"`
}

// entity for input for repository
type (
	CreateUserEntity struct {
		Username     string
		FirstName    string
		LastName     string
		PasswordHash string
	}

	UpdateUserEntity struct {
		UserID       string
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
		Username  string `json:"username"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Password  string `json:"password"`
	}

	LoginDTO struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	RefreshTokenDTO struct {
		RefreshToken string `json:"refresh_token"`
	}

	UpdateUserDTO struct {
		UserID    string
		Username  *string `json:"username"`
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Password  *string `json:"password"`
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
