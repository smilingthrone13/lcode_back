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
		TaskService         taskServ.Task
		TaskTemplateService taskTemplateServ.TaskTemplate
		TestCaseService     testCaseServ.TestCase
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

func (m *Manager) CreateProblem(ctx context.Context, dto domain.ProblemCreateDTO) (p domain.Problem, err error) {
	tx, err := m.transactionManager.NewTx(ctx, nil)
	if err != nil {
		return p, errors.Wrap(err, "ProblemManager Manager CreateProblem:")
	}
	ctx = context.WithValue(ctx, postgres.TxKey{}, tx)
	defer tx.Rollback(ctx)

	taskID, err := m.services.TaskService.Create(ctx, dto.Input.Task)
	if err != nil {
		return p, errors.Wrap(err, "ProblemManager Manager CreateProblem:")
	}

	for i := range dto.Input.TaskTemplates {
		err = m.services.TaskTemplateService.Create(ctx, taskID, dto.Input.TaskTemplates[i])
		if err != nil {
			return p, errors.Wrap(err, "ProblemManager Manager CreateProblem:")
		}
	}

	for i := range dto.Input.TestCases {
		err = m.services.TestCaseService.Create(ctx, taskID, dto.Input.TestCases[i])
		if err != nil {
			return p, errors.Wrap(err, "ProblemManager Manager CreateProblem:")
		}
	}

	p, err = m.FullProblemByTaskID(ctx, taskID)

	if err = tx.Commit(ctx); err != nil {
		return p, errors.Wrap(err, "ProblemManager Manager CreateProblem:")
	}

	return p, nil
}

func (m *Manager) UpdateProblemTask(
	ctx context.Context,
	dto domain.TaskUpdateDTO,
) (p domain.Problem, err error) {
	tx, err := m.transactionManager.NewTx(ctx, nil)
	if err != nil {
		return p, errors.Wrap(err, "ProblemManager Manager UpdateProblemTask:")
	}
	ctx = context.WithValue(ctx, postgres.TxKey{}, tx)
	defer tx.Rollback(ctx)

	err = m.services.TaskService.Update(ctx, dto.TaskID, dto.Input)
	if err != nil {
		return p, errors.Wrap(err, "ProblemManager Manager UpdateProblemTask:")
	}

	p, err = m.FullProblemByTaskID(ctx, dto.TaskID)
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
	dto domain.TaskTemplateCreateDTO,
) (p domain.Problem, err error) {
	tx, err := m.transactionManager.NewTx(ctx, nil)
	if err != nil {
		return p, errors.Wrap(err, "ProblemManager Manager CreateProblemTaskTemplate:")
	}
	ctx = context.WithValue(ctx, postgres.TxKey{}, tx)
	defer tx.Rollback(ctx)

	err = m.services.TaskTemplateService.Create(ctx, dto.TaskID, dto.Input)
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
	dto domain.TaskTemplateUpdateDTO,
) (p domain.Problem, err error) {
	tx, err := m.transactionManager.NewTx(ctx, nil)
	if err != nil {
		return p, errors.Wrap(err, "ProblemManager Manager UpdateProblemTaskTemplate:")
	}
	ctx = context.WithValue(ctx, postgres.TxKey{}, tx)
	defer tx.Rollback(ctx)

	err = m.services.TaskTemplateService.Update(ctx, dto.TemplateID, dto.Input)
	if err != nil {
		return p, errors.Wrap(err, "ProblemManager Manager UpdateProblemTaskTemplate:")
	}

	p, err = m.FullProblemByTaskID(ctx, dto.TaskID)
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

func (m *Manager) CreateProblemTestCase(
	ctx context.Context,
	dto domain.TestCaseCreateDTO,
) (p domain.Problem, err error) {
	tx, err := m.transactionManager.NewTx(ctx, nil)
	if err != nil {
		return p, errors.Wrap(err, "ProblemManager Manager CreateProblemTestCase:")
	}
	ctx = context.WithValue(ctx, postgres.TxKey{}, tx)
	defer tx.Rollback(ctx)

	err = m.services.TestCaseService.Create(ctx, dto.TaskID, dto.Input)
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
	dto domain.TestCaseUpdateDTO,
) (p domain.Problem, err error) {
	tx, err := m.transactionManager.NewTx(ctx, nil)
	if err != nil {
		return p, errors.Wrap(err, "ProblemManager Manager UpdateProblemTestCase:")
	}
	ctx = context.WithValue(ctx, postgres.TxKey{}, tx)
	defer tx.Rollback(ctx)

	err = m.services.TestCaseService.Update(ctx, dto.CaseID, dto.Input)
	if err != nil {
		return p, errors.Wrap(err, "ProblemManager Manager UpdateProblemTestCase:")
	}

	p, err = m.FullProblemByTaskID(ctx, dto.TaskID)
	if err != nil {
		return p, errors.Wrap(err, "ProblemManager Manager UpdateProblemTestCase:")
	}

	if err = tx.Commit(ctx); err != nil {
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

func (m *Manager) GetAvailableTaskAttributes(ctx context.Context) (domain.TaskAttributes, error) {
	ta, err := m.services.TaskService.GetAvailableAttributes(ctx)
	if err != nil {
		return ta, errors.Wrap(err, "ProblemManager Manager GetAvailableTaskAttributes:")
	}

	return ta, nil
}
