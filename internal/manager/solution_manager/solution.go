package solution_manager

import (
	"cmp"
	"context"
	"github.com/pkg/errors"
	"lcode/config"
	"lcode/internal/domain"
	"lcode/internal/service/solution"
	solutionResult "lcode/internal/service/solution_result"
	"lcode/pkg/postgres"
	"log"
	"log/slog"
	"slices"
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

		availableStatuses []domain.JudgeStatusInfo
	}
)

func New(
	cfg *config.Config,
	logger *slog.Logger,
	transactionManager *postgres.TransactionProvider,
	services *Services,
) *Manager {
	statuses, err := services.Judge.GetAvailableStatuses(context.Background())
	if err != nil {
		log.Fatal("can not access judge api:", err.Error())
	}

	m := &Manager{
		cfg:                cfg,
		logger:             logger,
		transactionManager: transactionManager,
		services:           services,
		solutionQueue:      newSolutionQueue(1000),
		workerCh:           make(chan workerItem, workersCount), // todo: config value
		availableStatuses:  statuses,
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
	baseCtx := context.Background()
	solUpdateStatus := domain.SolutionStatusCompleted
	sol := item.solution
	task := &item.task
	template := &item.template
	testCases := item.testCases

	solResults := make([]domain.SolutionResult, 0, len(testCases))
	srcCode := sol.Code + template.Wrapper

	for i := range testCases {
		data := domain.CreateJudgeSubmission{
			SourceCode:     srcCode,
			LanguageID:     sol.LanguageID,
			Stdin:          testCases[i].Input,
			ExpectedOutput: testCases[i].Output,
			CpuTimeLimit:   task.RuntimeLimit,
			MemoryLimit:    task.MemoryLimit,
		}

		info, err := m.createSubmission(baseCtx, data)
		if err != nil {
			solUpdateStatus = domain.SolutionStatusError

			m.logger.Error("can not create submission", slog.String("err", err.Error()))

			break
		}

		result := domain.SolutionResult{
			SolutionID:      sol.Id,
			TestCaseID:      testCases[i].ID,
			SubmissionToken: info.Token,
			Status:          info.Status,
			Runtime:         info.Time,
			Memory:          info.Memory,
			Stdout:          info.Stdout,
			Stderr:          info.Stderr,
		}

		solResults = append(solResults, result)

		if info.Status != domain.Accepted {
			solUpdateStatus = domain.SolutionStatusError
			break
		}
	}

	tx, err := m.transactionManager.NewTx(baseCtx, nil)
	if err != nil {
		m.logger.Error("can not create transaction", slog.String("err", err.Error()))

		return
	}
	ctx := context.WithValue(baseCtx, postgres.TxKey{}, tx)
	defer tx.Rollback(ctx)

	// если не получилось, то пропускаем и меняем статус у solution на error

	if len(solResults) != 0 {
		err = m.services.SolutionResult.CreateBatch(ctx, solResults...)
		if err != nil {
			solUpdateStatus = domain.SolutionStatusError

			m.logger.Error("can not create solution results", slog.String("err", err.Error()))
		}
	}

	var maxRuntimeSolResult domain.SolutionResult

	if len(solResults) != 0 {
		maxRuntimeSolResult = slices.MaxFunc(solResults, func(a, b domain.SolutionResult) int {
			return cmp.Compare(a.Runtime, b.Runtime)
		})
	}

	updateSolutionDTO := domain.UpdateSolutionDTO{
		ID:      sol.Id,
		Status:  &solUpdateStatus,
		Runtime: &maxRuntimeSolResult.Runtime,
		Memory:  &maxRuntimeSolResult.Memory,
	}

	_, err = m.services.Solution.Update(ctx, updateSolutionDTO)
	if err != nil {
		m.logger.Error("can not set status to solution", slog.String("err", err.Error()))

		return
	}

	if err = tx.Commit(ctx); err != nil {
		m.logger.Error("can not commit transaction", slog.String("err", err.Error()))

		return
	}
}

func (m *Manager) createSubmission(
	ctx context.Context,
	data domain.CreateJudgeSubmission,
) (info domain.JudgeSubmissionInfo, err error) {
	for {
		info, err = m.services.Judge.CreateSubmission(ctx, data)
		var queueIsFullError *domain.JudgeQueueIsFullError

		if errors.As(err, &queueIsFullError) {
			time.Sleep(time.Millisecond * 100)
			continue
		} else if err != nil {
			return domain.JudgeSubmissionInfo{}, errors.Wrap(err, "createSubmission solution manager")
		}

		return info, nil
	}
}

func (m *Manager) CreateSolution(
	ctx context.Context,
	dto domain.CreateSolutionDTO,
) (sol domain.Solution, err error) {
	tx, err := m.transactionManager.NewTx(ctx, nil)
	if err != nil {
		return domain.Solution{}, errors.Wrap(err, "CreateSolution solution manager")
	}
	ctx = context.WithValue(ctx, postgres.TxKey{}, tx)
	defer tx.Rollback(ctx)

	entity := domain.CreateSolutionEntity{
		TaskID:     dto.TaskID,
		LanguageID: dto.LanguageID,
		Code:       dto.Code,
		Status:     domain.SolutionStatusTesting,
		User:       dto.User,
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

func (m *Manager) GetAvailableSolutionStatuses() ([]domain.JudgeStatusInfo, error) {
	return m.availableStatuses, nil
}
