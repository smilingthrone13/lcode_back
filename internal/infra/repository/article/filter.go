package article

import (
	sql_query_maker "github.com/m-a-r-a-t/sql-query-maker"
	"lcode/config"
	"lcode/internal/domain"
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

func (f *filter) ConditionSearch(search string, searchCoefficient float32) *filter {
	if search != "" {
		f.Add("AND word_similarity(?, a.title) >= ?", search, searchCoefficient)
	}

	return f
}

func (f *filter) ConditionCategories(categories []string) *filter {
	if len(categories) > 0 {
		f.Add("AND ? && a.categories", categories)
	}

	return f
}

func (f *filter) AddCondition(p domain.ArticleParams) *filter {
	f.ConditionSearch(p.Filter.Search, f.conf.SearchCoefficient)
	f.ConditionCategories(p.Filter.Categories)

	return f
}

func (f *filter) WhereOptional(modifyFunc func()) *filter {
	f.SqlQueryMaker.WhereOptional(modifyFunc)

	return f
}

func (f *filter) SortByCreatedAt(t db.SortType) *filter {
	if t == db.ASC {
		f.Add("ORDER BY created_at, id")
	} else {
		f.Add("ORDER BY created_at DESC, id DESC")
	}

	return f
}

//func (f *filter) Sort(s domain.ArticleSort) *filter {
//	f.Add("ORDER BY")
//
//	switch s.ByTitle {
//	case db.ASC:
//		f.Add("title,")
//	case db.DESC:
//		f.Add("title DESC,")
//	}
//
//	if s.ByDate == db.DESC {
//		f.Add("created_at DESC")
//	} else {
//		f.Add("created_at")
//	}
//
//	return f
//}

func (f *filter) Clear() *filter {
	f.SqlQueryMaker.Clear()

	return f
}
