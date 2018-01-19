package main

import (
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
)

type Resource struct {
	queue     []*Event
	bandwidth float64

	lastTimeRequest float64
	mutex           sync.Mutex

	counter int64
	env     *Environment
}

func (r *Resource) Put(e *Event) {
	r.CounterAdd()

	r.update()
	r.queue = append(r.queue, e)
	sort.Sort(ByRemainingSize(r.queue))
	n := float64(atomic.LoadInt64(&r.counter))
	r.bandwidth *= n / (n + 1)
	if r.queue[0] == e {
		e.timeEnd = r.queue[0].remainingSize / r.bandwidth
	}
}

func (r *Resource) CounterAdd() {
	atomic.AddInt64(&r.counter, 1)
}

func (r *Resource) CounterMinus() {
	atomic.AddInt64(&r.counter, -1)
}

type Event struct {
	id        string
	size      float64
	timeStart float64
	timeEnd   interface{}

	remainingSize float64
	resource      interface{}
	callbacks     []func(*Event)
	worker        *Worker
	sender        string
	receiver      string
	listener      string

	recv bool
	send bool

	twinEvent *Event
}

type ImmuteEvent struct {
	Event
}

func (r *Resource) update() {
	//Delete remaining size

	for index := range r.queue {
		r.queue[index].remainingSize -= (r.env.currentTime - r.lastTimeRequest) * r.bandwidth
	}
}

type ByTime []*Event

func (s ByTime) Len() int {
	return len(s)
}
func (s ByTime) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByTime) Less(i, j int) bool {
	if s[i].timeEnd == nil && s[j].timeEnd == nil {
		return s[i].size < s[j].size
	} else if s[i].timeEnd == nil {
		return false
	} else if s[j].timeEnd == nil {
		return true
	} else {
		return s[i].timeEnd.(float64) < s[j].timeEnd.(float64)
	}
}

func (e *Event) String() string {
	return fmt.Sprintf("%v", e.timeEnd)
}

type ByRemainingSize []*Event

func (b ByRemainingSize) Len() int {
	return len(b)
}
func (b ByRemainingSize) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
func (b ByRemainingSize) Less(i, j int) bool {
	return b[i].remainingSize < b[j].remainingSize
}
