package work_queue

import (
	"errors"
	"lib/lru_cache"
	"runtime/debug"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	ErrorWorkQueueClosed = errors.New("err work queue closed")
	ErrorWorkQueueFulled = errors.New("err work queue fulled")
)

type WorkQueue struct {
	sync.RWMutex
	sync.WaitGroup
	pipe     int
	cache    int
	workChan []chan func()
	isRun    bool
	prefix   string
}

func CaptureException() {
	if err := recover(); err != nil {
		logrus.Error("Recovered in err", err, string(debug.Stack()))
	}
}

func NewWorkQueue(pipe, cache int, prefix string) *WorkQueue {
	ret := &WorkQueue{
		pipe:   pipe,
		cache:  cache,
		prefix: prefix,
	}
	ret.Run()
	return ret
}

func (r *WorkQueue) Run() {
	r.Lock()
	defer r.Unlock()
	logrus.Infof("work queue:%s start with pipe:%d cache:%d", r.prefix, r.pipe, r.cache)
	r.workChan = make([]chan func(), r.pipe)
	for i := 0; i < r.pipe; i++ {
		r.workChan[i] = make(chan func(), r.cache)
		go r.process(i)
	}
	r.isRun = true
}

func (r *WorkQueue) SendTask(fn func()) error {
	shardId := int(time.Now().UnixNano()) % r.pipe
	if !r.isRun {
		logrus.Errorf("err work_queue:%s closed", r.prefix)
		return ErrorWorkQueueClosed
	}
	select {
	case r.workChan[shardId] <- fn:
		return nil
	default:
		lru_cache.FreqCall("workqueue:"+r.prefix, time.Second, func() {
			logrus.Errorf("workqueue:%s shardId:%d fulled", r.prefix, shardId)
		})
		return ErrorWorkQueueFulled
	}
}

func (r *WorkQueue) process(i int) {
	r.Add(1)
	defer r.Done()
	workChan := <-r.workChan
	defer func() { r.workChan[i] = nil }()
	for {
		task, ok := <-workChan
		if !ok {
			return
		}
		fn := func() {
			CaptureException()
			task()
		}()
		fn()
	}
}

func (r *WorkQueue) Stop() {
	r.Lock()
	defer r.Unlock()
	r.is_run = false
	logrus.Infof("work_queue:%s stopping", r.prefix)
	for _, workChan := range r.workChan {
		close(workChan)
	}
	r.Wait()
	logrus.Infof("work queue:%s stopped", r.prefix)
}
