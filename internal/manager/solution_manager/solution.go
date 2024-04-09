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
	"time"
)

const (
	workersCount = 8
)

type (
	Services struct {
		ProblemManager ProblemManager
		Solution       solution.Solution
		SolutionResult solutionResult.SolutionResult
		Judge          Judge
	}

	Manager struct {
		cfg                *config.Config
		logger             *slog.Logger
		transactionManager *postgres.TransactionProvider
		services           *Services
		solutionQueue      *solutionQueue
		workerCh           chan workerItem
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
		workerCh:           make(chan workerItem, workersCount), // todo: config value
	}

	go m.runWorkerManager()

	for _ = range workersCount {
		go func() {
			for item := range m.workerCh {
				m.solutionWorker(item)
			}
		}()
	}

	return m
}

func (m *Manager) runWorkerManager() {
	ticker := time.NewTicker(time.Millisecond * 10)
	defer ticker.Stop()

	for {
		<-ticker.C

		sol, ok := m.solutionQueue.PopFront()
		if !ok {
			continue
		}

		ctx := context.Background()

		problem, err := m.services.ProblemManager.FullProblemByTaskID(ctx, sol.TaskID)
		if err != nil {
			m.logger.Error("can not find problem by task_id", slog.String("err", err.Error()))
			continue
		}

		var tmpl *domain.TaskTemplate

		for i := range problem.TaskTemplates {
			if problem.TaskTemplates[i].LanguageID == sol.LanguageID {
				tmpl = &problem.TaskTemplates[i]
				break
			}
		}

		if tmpl == nil {
			m.logger.Error(
				"template for user solution was not found in the task",
				slog.String("solution_id", sol.Id),
			)

			s := domain.SolutionStatusError
			_, err = m.services.Solution.Update(ctx, domain.UpdateSolutionDTO{
				ID:     sol.Id,
				Status: &s,
			})

			if err != nil {
				m.logger.Error("can not update solution status to error", slog.String("err", err.Error()))
			}

			continue
		}

		item := workerItem{
			solution:  sol,
			task:      problem.Task,
			template:  *tmpl,
			testCases: problem.TestCases,
		}

		m.workerCh <- item
	}

}

func (m *Manager) solutionWorker(item workerItem) {
	//ctx := context.Background()

	//task := &item.task
	//template := &item.template
	//testCases := item.testCases
	//
	//data := domain.CreateJudgeSubmission{
	//	SourceCode: solution.Code,
	//	LanguageID: solution.LanguageID,
	//	Stdin: item.task.
	//}
	//
	//m.services.Judge.CreateSubmission(ctx,)
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
		Status:     domain.SolutionStatusTesting,
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
