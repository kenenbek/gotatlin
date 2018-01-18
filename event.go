package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Resource struct {
	queue     []*Event
	bandwidth float64

	lastTimeRequest float64
	mutex           sync.Mutex

	counter int64
}

func (r *Resource) Put(e *Event) float64{
	
	return
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
	timeEnd   float64

	remainingSize float64
	resource      interface{}
	callbacks     []func(*Event)
	worker *Worker
	sender string
	receiver string
	listener string

	recv bool
	send bool
}

type ImmuteEvent struct {
	Event
}

func (e *Event) update(deltaT float64, env *Environment) {
	switch e.resource.(type) {
	case Link, *Link:
		if e.timeStart < env.currentTime {
			e.remainingSize -= deltaT * e.resource.(*Link).bandwidth
		}
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
	return s[i].timeEnd < s[j].timeEnd
}

func (e *Event) String() string {
	return fmt.Sprintf("%v", e.timeEnd)
}
