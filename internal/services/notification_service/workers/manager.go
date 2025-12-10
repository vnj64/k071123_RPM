package workers

import "context"

type Consumer interface {
	Start(ctx context.Context) error
}

type WorkerManager struct {
	consumers []Consumer
}

func NewWorkerManager() *WorkerManager {
	return &WorkerManager{}
}

func (m *WorkerManager) Register(consumer Consumer) {
	m.consumers = append(m.consumers, consumer)
}

func (m *WorkerManager) StartAll(ctx context.Context) error {
	for _, c := range m.consumers {
		if err := c.Start(ctx); err != nil {
			return err
		}
	}
	return nil
}
