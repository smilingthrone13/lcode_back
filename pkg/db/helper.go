package db

type SortType string

const (
	ASC  SortType = "asc"
	DESC SortType = "desc"
)

func GetLetterGreaterOrLessBySortType(t SortType) string {
	if t == DESC {
		return "<"
	} else {
		return ">"
	}
}
