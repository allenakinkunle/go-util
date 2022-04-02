package asyncer

import (
	"github.com/hibiken/asynq"
	"time"
)

type IAsyncer interface {
	EnqueueTask(taskName, taskID string, payload []byte) error
	ScheduleTask(taskName, taskID string, payload []byte, in time.Duration) error
}

type Asyncer struct {
	client *asynq.Client
}

func NewAsyncer(host, port string, db int) *Asyncer {
	return &Asyncer{
		client: asynq.NewClient(asynq.RedisClientOpt{
			Addr: host + ":" + port,
			DB:   db,
		}),
	}
}

func (a Asyncer) EnqueueTask(taskName, taskID string, payload []byte) error {
	task := asynq.NewTask(taskName, payload)
	_, err := a.client.Enqueue(task, asynq.TaskID(taskID))
	return err
}

func (a Asyncer) ScheduleTask(taskName, taskID string, payload []byte, in time.Duration) error {
	task := asynq.NewTask(taskName, payload)
	_, err := a.client.Enqueue(task, asynq.TaskID(taskID), asynq.ProcessIn(in))
	return err
}
