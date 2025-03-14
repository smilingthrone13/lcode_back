package solution_manager

import (
	"container/list"
	"github.com/pkg/errors"
	"lcode/internal/domain"
	"sync"
)

func newSolutionQueue(maxSize int) *solutionQueue {
	return &solutionQueue{
		maxSize: maxSize,
		queue:   list.New(),
	}
}

type solutionQueue struct {
	maxSize int
	mu      sync.Mutex
	queue   *list.List
}

func (q *solutionQueue) PushBack(s domain.Solution) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.queue.Len() >= q.maxSize {
		return errors.New("queue is full")
	}

	q.queue.PushBack(s)

	return nil
}

func (q *solutionQueue) PopFront() (domain.Solution, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.queue.Len() == 0 {
		return domain.Solution{}, false
	}

	s := q.queue.Front()

	q.queue.Remove(s)

	return s.Value.(domain.Solution), true
}

func (q *solutionQueue) PopFrontAll() []domain.Solution {
	q.mu.Lock()
	defer q.mu.Unlock()

	solutions := make([]domain.Solution, 0, q.queue.Len())

	var e *list.Element
	for q.queue.Len() > 0 {
		e = q.queue.Front()

		solutions = append(solutions, e.Value.(domain.Solution))

		q.queue.Remove(e)
	}

	return solutions
}
