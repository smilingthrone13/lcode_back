package solution_manager

import "lcode/internal/domain"

type workerItem struct {
	solution  domain.Solution
	task      domain.Task
	template  domain.TaskTemplate
	testCases []domain.TestCase
}
