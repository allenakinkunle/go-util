package asyncer

import "time"

type mockAsyncer struct {
}

func NewMockAsyncer() *mockAsyncer {
	return &mockAsyncer{}
}

func (m mockAsyncer) EnqueueTask(taskName, taskID string, payload []byte) error {
	return nil
}

func (m mockAsyncer) ScheduleTask(taskName, taskID string, payload []byte, in time.Duration) error {
	return nil
}
