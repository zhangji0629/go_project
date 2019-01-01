package work_queue_test

import (
	"lib/work_queue"
	"testing"
)

func TestQueue1(t *testing.T) {
	worker1 := work_queue.NewWorkQueue(10, 20, "test1")
	worker2 := work_queue.NewWorkQueue(10, 20, "test2")

	for i := 1; i <= 100; i++ {
		j := i
		worker1.SendTask(func() {
			t.Log("hehehda:", j)
		})
	}

	for i := 1; i <= 100; i++ {
		j := i
		worker2.SendTask(func() {
			t.Log("lalala:", j)
		})
	}

	worker1.Stop()
	worker2.Stop()
}
