package domain

type (
	User struct {
		ID        string
		Login     string
		FirstName string
		LastName  string
		IsAdmin   bool
	}
)
