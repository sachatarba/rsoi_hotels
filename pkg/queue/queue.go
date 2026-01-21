package queue

import (
	"log/slog"
	"time"
)

type Task func() error

type Queue struct {
	tasks  chan Task
	logger *slog.Logger
}

func New(logger *slog.Logger) *Queue {
	return &Queue{
		tasks:  make(chan Task, 100),
		logger: logger,
	}
}

func (q *Queue) Add(task Task) {
	q.tasks <- task
}

func (q *Queue) StartWorker() {
	go func() {
		for task := range q.tasks {
			go q.processTask(task)
		}
	}()
}

func (q *Queue) processTask(task Task) {
	for {
		err := task()
		if err == nil {
			q.logger.Info("Background task completed successfully")
			return
		}

		q.logger.Error("Background task failed, will retry in 10 seconds", "error", err)
		time.Sleep(10 * time.Second)
	}
}
