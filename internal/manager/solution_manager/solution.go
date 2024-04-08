package solution_manager

import (
	"context"
	"github.com/pkg/errors"
	"lcode/config"
	"lcode/internal/domain"
	"lcode/internal/service/solution"
	solutionResult "lcode/internal/service/solution_result"
	"lcode/pkg/postgres"
	"log/slog"
)

const (
	workersCount = 8
)

type (
	Services struct {
		Solution       solution.Solution
		SolutionResult solutionResult.SolutionResult
	}

	Manager struct {
		cfg                *config.Config
		logger             *slog.Logger
		transactionManager *postgres.TransactionProvider
		services           *Services
		solutionQueue      *solutionQueue
		workerCh           chan domain.Solution
	}
)

func New(
	cfg *config.Config,
	logger *slog.Logger,
	transactionManager *postgres.TransactionProvider,
	services *Services,
) *Manager {

	m := &Manager{
		cfg:                cfg,
		logger:             logger,
		transactionManager: transactionManager,
		services:           services,
		solutionQueue:      newSolutionQueue(1000),
		workerCh:           make(chan domain.Solution, workersCount), // todo: config value
	}

	go m.runWorkerManager()

	for _ = range workersCount {
		go func() {
			for sol := range m.workerCh {
				m.solutionWorker(sol)
			}
		}()
	}

	return m
}

func (m *Manager) runWorkerManager() {
	// todo
}

func (m *Manager) solutionWorker(sol domain.Solution) {

}

func (m *Manager) CreateSolution(ctx context.Context, dto domain.CreateSolutionDTO) (sol domain.Solution, err error) {
	tx, err := m.transactionManager.NewTx(ctx, nil)
	if err != nil {
		return domain.Solution{}, errors.Wrap(err, "CreateSolution solution manager")
	}
	ctx = context.WithValue(ctx, postgres.TxKey{}, tx)
	defer tx.Rollback(ctx)

	entity := domain.CreateSolutionEntity{
		UserID:     dto.UserID,
		TaskID:     dto.TaskID,
		LanguageID: dto.LanguageID,
		Code:       dto.Code,
		Status:     domain.Testing,
	}

	sol, err = m.services.Solution.Create(ctx, entity)
	if err != nil {
		return domain.Solution{}, errors.Wrap(err, "CreateSolution solution manager")
	}

	err = m.solutionQueue.PushBack(sol)
	if err != nil {
		return domain.Solution{}, errors.Wrap(err, "CreateSolution solution manager")
	}

	if err = tx.Commit(ctx); err != nil {
		return domain.Solution{}, errors.Wrap(err, "CreateSolution solution manager")
	}

	return sol, nil
}
