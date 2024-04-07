package solution

import "context"

type (
	Solution interface {
	}

	SolutionRepo interface {
		Create(ctx context.Context)
		ChangeStatus(ctx context.Context)
	}
)
