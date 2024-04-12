package comment

import (
	sql_query_maker "github.com/m-a-r-a-t/sql-query-maker"
	"lcode/config"
	"lcode/pkg/db"
)

type filter struct {
	conf *config.Config
	*sql_query_maker.SqlQueryMaker
}

func newFilter(conf *config.Config, argsCount int) *filter {
	return &filter{
		conf:          conf,
		SqlQueryMaker: sql_query_maker.NewQueryMaker(argsCount),
	}
}

func (f *filter) Add(query string, args ...interface{}) *filter {
	f.SqlQueryMaker.Add(query, args...)

	return f
}

func (f *filter) SortByCreatedAt(t db.SortType) *filter {
	if t == db.ASC {
		f.Add("ORDER BY c.created_at, id")
	} else {
		f.Add("ORDER BY c.created_at DESC, id DESC")
	}

	return f
}

func (f *filter) Clear() *filter {
	f.SqlQueryMaker.Clear()

	return f
}
