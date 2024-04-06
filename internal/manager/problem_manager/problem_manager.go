package problem_manager

import (
	"context"
	"github.com/pkg/errors"
	"lcode/config"
	"lcode/internal/domain"
	taskServ "lcode/internal/service/task"
	taskTemplateServ "lcode/internal/service/task_template"
	testCaseServ "lcode/internal/service/test_case"
	"lcode/pkg/postgres"
	"log/slog"
)

type (
	Services struct {
		TaskService         *taskServ.Service
		TaskTemplateService *taskTemplateServ.Service
		TestCaseService     *testCaseServ.Service
	}

	Manager struct {
		cfg                *config.Config
		logger             *slog.Logger
		transactionManager *postgres.TransactionProvider
		services           *Services
	}
)

func New(
	cfg *config.Config,
	logger *slog.Logger,
	transactionManager *postgres.TransactionProvider,
	services *Services,
) *Manager {
	return &Manager{
		cfg:                cfg,
		logger:             logger,
		transactionManager: transactionManager,
		services:           services,
	}
}

func (m *Manager) CreateProblem(ctx context.Context, dto domain.ProblemCreateInput) (p domain.Problem, err error) {
	tx, err := m.transactionManager.NewTx(ctx, nil)
	if err != nil {
		return p, errors.Wrap(err, "ProblemManager Manager CreateProblem:")
	}
	ctx = context.WithValue(ctx, postgres.TxKey{}, tx)
	defer tx.Rollback(ctx)

	task, err := m.services.TaskService.Create(ctx, dto.Task)
	if err != nil {
		return p, errors.Wrap(err, "ProblemManager Manager CreateProblem:")
	}

	templates := make([]domain.TaskTemplate, 0, len(dto.TaskTemplates))
	for i := range dto.TaskTemplates {
		template, err := m.services.TaskTemplateService.Create(ctx, dto.TaskTemplates[i])
		if err != nil {
			return p, errors.Wrap(err, "ProblemManager Manager CreateProblem:")
		}

		templates = append(templates, template)
	}

	testCases := make([]domain.TestCase, 0, len(dto.TestCases))
	for i := range dto.TestCases {
		testCase, err := m.services.TestCaseService.Create(ctx, dto.TestCases[i])
		if err != nil {
			return p, errors.Wrap(err, "ProblemManager Manager CreateProblem:")
		}

		testCases = append(testCases, testCase)
	}

	if err = tx.Commit(ctx); err != nil {
		return p, errors.Wrap(err, "ProblemManager Manager CreateProblem:")
	}

	p = domain.Problem{
		Task:          task,
		TaskTemplates: templates,
		TestCases:     testCases,
	}

	return p, nil
}

func (m *Manager) UpdateProblemTask(
	ctx context.Context,
	taskID string,
	dto domain.TaskUpdateInput,
) (p domain.Problem, err error) {
	tx, err := m.transactionManager.NewTx(ctx, nil)
	if err != nil {
		return p, errors.Wrap(err, "ProblemManager Manager UpdateProblemTask:")
	}
	ctx = context.WithValue(ctx, postgres.TxKey{}, tx)
	defer tx.Rollback(ctx)

	_, err = m.services.TaskService.Update(ctx, taskID, dto)
	if err != nil {
		return p, errors.Wrap(err, "ProblemManager Manager UpdateProblemTask:")
	}

	p, err = m.FullProblemByTaskID(ctx, taskID)
	if err != nil {
		return p, errors.Wrap(err, "ProblemManager Manager UpdateProblemTask:")
	}

	if err = tx.Commit(ctx); err != nil {
		return p, errors.Wrap(err, "ProblemManager Manager UpdateProblemTask:")
	}

	return p, nil
}

func (m *Manager) DeleteProblem(ctx context.Context, taskID string) (err error) {
	tx, err := m.transactionManager.NewTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "ProblemManager Manager DeleteProblem:")
	}
	ctx = context.WithValue(ctx, postgres.TxKey{}, tx)
	defer tx.Rollback(ctx)

	err = m.services.TaskService.Delete(ctx, taskID)
	if err != nil {
		return errors.Wrap(err, "ProblemManager Manager DeleteProblem:")
	}

	if err = tx.Commit(ctx); err != nil {
		return errors.Wrap(err, "ProblemManager Manager DeleteProblem:")
	}

	return nil
}

func (m *Manager) CreateProblemTaskTemplate(
	ctx context.Context,
	dto domain.TaskTemplateCreateInput,
) (p domain.Problem, err error) {
	tx, err := m.transactionManager.NewTx(ctx, nil)
	if err != nil {
		return p, errors.Wrap(err, "ProblemManager Manager CreateProblemTaskTemplate:")
	}
	ctx = context.WithValue(ctx, postgres.TxKey{}, tx)
	defer tx.Rollback(ctx)

	_, err = m.services.TaskTemplateService.Create(ctx, dto)
	if err != nil {
		return p, errors.Wrap(err, "ProblemManager Manager CreateProblemTaskTemplate:")
	}

	p, err = m.FullProblemByTaskID(ctx, dto.TaskID)
	if err != nil {
		return p, errors.Wrap(err, "ProblemManager Manager CreateProblemTaskTemplate:")
	}

	if err = tx.Commit(ctx); err != nil {
		return p, errors.Wrap(err, "ProblemManager Manager CreateProblemTaskTemplate:")
	}

	return p, nil
}

func (m *Manager) UpdateProblemTaskTemplate(
	ctx context.Context,
	templateID string,
	dto domain.TaskTemplateUpdateInput,
) (p domain.Problem, err error) {
	tx, err := m.transactionManager.NewTx(ctx, nil)
	if err != nil {
		return p, errors.Wrap(err, "ProblemManager Manager UpdateProblemTaskTemplate:")
	}
	ctx = context.WithValue(ctx, postgres.TxKey{}, tx)
	defer tx.Rollback(ctx)

	template, err := m.services.TaskTemplateService.Update(ctx, templateID, dto)
	if err != nil {
		return p, errors.Wrap(err, "ProblemManager Manager UpdateProblemTaskTemplate:")
	}

	p, err = m.FullProblemByTaskID(ctx, template.TaskID)
	if err != nil {
		return p, errors.Wrap(err, "ProblemManager Manager UpdateProblemTaskTemplate:")
	}

	if err = tx.Commit(ctx); err != nil {
		return p, errors.Wrap(err, "ProblemManager Manager UpdateProblemTaskTemplate:")
	}

	return p, nil
}

func (m *Manager) DeleteProblemTaskTemplate(ctx context.Context, templateID string) error {
	tx, err := m.transactionManager.NewTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "ProblemManager Manager DeleteProblemTaskTemplate:")
	}
	ctx = context.WithValue(ctx, postgres.TxKey{}, tx)
	defer tx.Rollback(ctx)

	err = m.services.TaskTemplateService.Delete(ctx, templateID)
	if err != nil {
		return errors.Wrap(err, "ProblemManager Manager DeleteProblemTaskTemplate:")
	}

	if err = tx.Commit(ctx); err != nil {
		return errors.Wrap(err, "ProblemManager Manager DeleteProblemTaskTemplate:")
	}

	return nil
}

func (m *Manager) CreateProblemTestCase(ctx context.Context, dto domain.TestCaseCreateInput) (p domain.Problem, err error) {
	tx, err := m.transactionManager.NewTx(ctx, nil)
	if err != nil {
		return p, errors.Wrap(err, "ProblemManager Manager CreateProblemTestCase:")
	}
	ctx = context.WithValue(ctx, postgres.TxKey{}, tx)
	defer tx.Rollback(ctx)

	_, err = m.services.TestCaseService.Create(ctx, dto)
	if err != nil {
		return p, errors.Wrap(err, "ProblemManager Manager CreateProblemTestCase:")
	}

	p, err = m.FullProblemByTaskID(ctx, dto.TaskID)
	if err != nil {
		return p, errors.Wrap(err, "ProblemManager Manager CreateProblemTestCase:")
	}

	if err = tx.Commit(ctx); err != nil {
		return p, errors.Wrap(err, "ProblemManager Manager CreateProblemTestCase:")
	}

	return p, nil
}

func (m *Manager) UpdateProblemTestCase(ctx context.Context,
	caseID string,
	dto domain.TestCaseUpdateInput,
) (p domain.Problem, err error) {
	tx, err := m.transactionManager.NewTx(ctx, nil)
	if err != nil {
		return p, errors.Wrap(err, "ProblemManager Manager UpdateProblemTestCase:")
	}
	ctx = context.WithValue(ctx, postgres.TxKey{}, tx)
	defer tx.Rollback(ctx)

	testCase, err := m.services.TestCaseService.Update(ctx, caseID, dto)
	if err != nil {
		return p, errors.Wrap(err, "ProblemManager Manager UpdateProblemTestCase:")
	}

	p, err = m.FullProblemByTaskID(ctx, testCase.TaskID)
	if err != nil {
		return p, errors.Wrap(err, "ProblemManager Manager UpdateProblemTestCase:")
	}

	return p, nil
}

func (m *Manager) DeleteProblemTestCase(ctx context.Context, caseID string) error {
	tx, err := m.transactionManager.NewTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "ProblemManager Manager DeleteProblemTestCase:")
	}
	ctx = context.WithValue(ctx, postgres.TxKey{}, tx)
	defer tx.Rollback(ctx)

	err = m.services.TestCaseService.Delete(ctx, caseID)
	if err != nil {
		return errors.Wrap(err, "ProblemManager Manager DeleteProblemTestCase:")
	}

	if err = tx.Commit(ctx); err != nil {
		return errors.Wrap(err, "ProblemManager Manager DeleteProblemTestCase:")
	}

	return nil
}

func (m *Manager) FullProblemByTaskID(ctx context.Context, taskID string) (p domain.Problem, err error) {
	tx, err := m.transactionManager.NewTx(ctx, nil)
	if err != nil {
		return p, errors.Wrap(err, "ProblemManager Manager FullProblemByTaskID:")
	}
	ctx = context.WithValue(ctx, postgres.TxKey{}, tx)
	defer tx.Rollback(ctx)

	task, err := m.services.TaskService.GetByID(ctx, taskID)
	if err != nil {
		return p, errors.Wrap(err, "ProblemManager Manager FullProblemByTaskID:")
	}

	taskTemplates, err := m.services.TaskTemplateService.GetAllByTaskID(ctx, taskID)
	if err != nil {
		return p, errors.Wrap(err, "ProblemManager Manager FullProblemByTaskID:")
	}

	testCases, err := m.services.TestCaseService.GetAllByTaskID(ctx, taskID)
	if err != nil {
		return p, errors.Wrap(err, "ProblemManager Manager FullProblemByTaskID:")
	}

	p = domain.Problem{
		Task:          task,
		TaskTemplates: taskTemplates,
		TestCases:     testCases,
	}

	return p, nil
}

func (m *Manager) TaskListByParams(ctx context.Context, dto domain.TaskParams) (tl domain.TaskList, err error) {
	tl, err = m.services.TaskService.GetAllByParams(ctx, dto)
	if err != nil {
		return tl, errors.Wrap(err, "ProblemManager Manager TaskListByParams:")
	}

	return tl, nil
}
