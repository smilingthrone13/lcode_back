package solution_result

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	sql_query_maker "github.com/m-a-r-a-t/sql-query-maker"
	"github.com/pkg/errors"
	"lcode/internal/domain"
	"lcode/pkg/postgres"
)

func New(db *postgres.DbManager) *Repository {
	return &Repository{db: db}
}

type Repository struct {
	db *postgres.DbManager
}

func (r *Repository) CreateBatch(ctx context.Context, results ...domain.SolutionResult) error {
	sq := sql_query_maker.NewQueryMaker(4)

	sq.Add(`INSERT INTO solution_result 
    			  (solution_id, test_case_id, submission_token, status, runtime, memory, stdout, stderr)`,
	)

	for i := range results {
		sq.Values(
			results[i].SolutionID,
			results[i].TestCaseID,
			results[i].SubmissionToken,
			results[i].Status,
			results[i].Runtime,
			results[i].Memory,
			results[i].Stdout,
			results[i].Stderr,
		)
	}

	query, args := sq.Make()
	_, err := r.db.TxOrDB(ctx).Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "CreateBatch solution_result repo")
	}

	return nil
}

func (r *Repository) ResultsBySolutionID(ctx context.Context, solutionID string) ([]domain.SolutionResult, error) {
	sq := sql_query_maker.NewQueryMaker(1)

	results := []domain.SolutionResult{}

	sq.Add(`
			SELECT solution_id, test_case_id, submission_token, 
			       status, runtime, memory, stdout, stderr 
			FROM solution_result
			WHERE solution_id = ?`,
		solutionID,
	)

	query, args := sq.Make()

	err := pgxscan.Select(ctx, r.db.TxOrDB(ctx), &results, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "ResultsBySolutionID solution_result repo")
	}

	return results, nil
}
