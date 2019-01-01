package time_trace

import (
	"math"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type (
	tracevalue struct {
		count      int64
		total_time float64
		max_time   float64
		min_time   float64
	}

	Trace struct {
		prefix     string
		func_names *sync.Map
		check_time int64
	}

	TraceItem struct {
		t     *Trace
		begin float64
		name  string
	}
)

func UnixTimeFloat() float64 {
	return float64(time.Now().UnixNano()) / 1000000000.0
}

func (t *tracevalue) Add(count int64, total_time float64) {
	t.total_time += total_time
	t.count += count
	if t.max_time < total_time {
		t.max_time = total_time
	}
	if t.min_time > total_time {
		t.min_time = total_time
	}
}

func (t *tracevalue) traceinfo() (float64, float64, float64) {
	if t.count <= 0 {
		return 0, 0, 0
	}
	return t.max_time, t.min_time, t.total_time / float64(t.count)
}

func NewTrace(check_time int64, prefix string) *Trace {
	if check_time <= 0 {
		check_time = default_check_time
	}
	t := &Trace{
		check_time: check_time,
		func_names: &sync.Map{},
		prefix:     prefix,
	}
	go t.tickTraceInfo()
	return t
}

func (t *Trace) tickTraceInfo() {
	tick := time.NewTicker(time.Second * time.Duration(t.check_time))
	for {
		if _, ok := <-tick.C; ok {
			old := t.func_names
			t.func_names = &sync.Map{}
			t.traceInfo(old)
		} else {
			return
		}
	}
}

func (t *Trace) traceInfo(info *sync.Map) {
	info.Range(func(name, funcs interface{}) bool {
		item := funcs.(*tracevalue)
		max_time, min_time, avg_time := item.traceinfo()
		max_time *= 1000
		min_time *= 1000
		avg_time *= 1000
		//fmt.Printf("%s||%s||count:%d||max_time: %.3f ms||min_time: %.3f ms||avg_time: %.3f ms\n", t.prefix, name, item.count, max_time, min_time, avg_time)
		logrus.Infof("%s||%s||count:%d||max_time: %.3f ms||min_time: %.3f ms||avg_time: %.3f ms", t.prefix, name, item.count, max_time, min_time, avg_time)
		return true
	})
}

func (t *Trace) Add(name string, delta float64) {
	var item *tracevalue
	if v, ok := t.func_names.Load(name); ok {
		item = v.(*tracevalue)
	} else {
		item = &tracevalue{
			max_time:   -math.MaxFloat64,
			min_time:   math.MaxFloat64,
			count:      0,
			total_time: 0,
		}
	}
	item.Add(1, delta)
	t.func_names.Store(name, item)
}

func (t *Trace) Begin(name string) *TraceItem {
	begin := UnixTimeFloat()
	return &TraceItem{
		name:  name,
		t:     t,
		begin: begin,
	}
}

func (i *TraceItem) End() {
	delta := UnixTimeFloat() - i.begin
	if i.t != nil {
		i.t.Add(i.name, delta)
	}
}

var (
	default_check_time int64 = 1
	GPacketTimeTrace         = NewTrace(default_check_time, "PacketTimeTrace")
	GRedisTimeTrace          = NewTrace(default_check_time, "RedisTimeTrace")
	GMysqlTimeTrace          = NewTrace(default_check_time, "MysqlTimeTrace")
	GHttpTimeTrace           = NewTrace(default_check_time, "HttpTimeTrace")
)
