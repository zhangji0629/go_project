package work_queue_test

import (
	"fmt"
	"lib/work_queue"
	"testing"
)

func TestQueue1(t *testing.T) {
	worker1 := work_queue.NewWorkQueue(10, 20, "test1")
	worker2 := work_queue.NewWorkQueue(10, 20, "test2")

	for i := 1; i <= 100; i++ {
		worker1.SendTask(func() {
			fmt.Println("hehehda", i)
		})
	}

	for i := 1; i <= 100; i++ {
		worker2.SendTask(func() {
			fmt.Println("hehehda", i)
		})
	}
}
