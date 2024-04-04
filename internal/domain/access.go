package domain

type (
	User struct {
		ID        string
		Login     string
		FirstName string
		LastName  string
		isAdmin   bool
	}
)
