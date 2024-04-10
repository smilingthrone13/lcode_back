package task

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
		f.Add("AND word_similarity(?, t.name) >= ?", search, searchCoefficient)
	}

	return f
}

func (f *filter) ConditionCategories(categories []string) *filter {
	if len(categories) > 0 {
		f.Add("AND t.category = ANY(?)", categories)
	}

	return f
}

func (f *filter) ConditionDifficulties(difficulties []string) *filter {
	if len(difficulties) > 0 {
		f.Add("AND t.difficulty = ANY(?)", difficulties)
	}

	return f
}

func (f *filter) AddCondition(p domain.TaskParams) *filter {
	f.ConditionSearch(p.Filter.Search, f.conf.SearchCoefficient)
	f.ConditionCategories(p.Filter.Categories)
	f.ConditionDifficulties(p.Filter.Difficulties)

	return f
}

func (f *filter) WhereOptional(modifyFunc func()) *filter {
	f.SqlQueryMaker.WhereOptional(modifyFunc)

	return f
}

func (f *filter) SortByNumber(t db.SortType) *filter {
	if t == db.DESC {
		f.Add("ORDER BY number DESC, id DESC")
	} else {
		f.Add("ORDER BY number, id")
	}

	return f
}

func (f *filter) Clear() *filter {
	f.SqlQueryMaker.Clear()

	return f
}
