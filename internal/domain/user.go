package domain

const (
	UserCtxKey = "user"
	DtoCtxKey  = "dto"
)

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	PasswordHash string `json:"-"`
	IsAdmin      bool   `json:"is_admin"`
}

type CreateUserDTO struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

type CreateUserEntity struct {
	Username     string
	FirstName    string
	LastName     string
	PasswordHash string
}

type LoginDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token"`
}

type ChangeUserAdminPermissionDTO struct {
	UserID  string `json:"user_id"`
	IsAdmin bool   `json:"is_admin"`
}
