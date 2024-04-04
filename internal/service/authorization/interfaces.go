package authorization

import (
	"lcode/internal/domain"
)

type Authorization interface {
	UserByID(id string) (user domain.User, ok bool)
	Users() ([]domain.User, error)
}
