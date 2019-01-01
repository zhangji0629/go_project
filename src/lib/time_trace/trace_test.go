package time_trace_test

import (
	"fmt"
	"lib/time_trace"
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func func2() {
	defer time_trace.GPacketTimeTrace.Begin("func2").End()
	xx := rand.Int31n(50) + 50
	time.Sleep(time.Millisecond * time.Duration(xx))
	fmt.Println(time.Now().Unix(), ": ", rand.Int63n(100), xx)
}

func func1() {
	tick := time.NewTicker(time.Millisecond * time.Duration(100))
	for {
		if _, ok := <-tick.C; ok {
			go func2()
		}
	}
}

func TestTrace_Count(t *testing.T) {
	func1()
}
