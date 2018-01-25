package main

import (
	"sort"
	"sync"
	"sync/atomic"
)

type Resource struct {
	queue     []*TransferEvent
	bandwidth float64

	lastTimeRequest float64
	mutex           sync.Mutex

	counter int64
	env     *Environment
}

func (r *Resource) Put(e *TransferEvent) {
	r.CounterAdd()

	r.queue = append(r.queue, e)
	sort.Sort(ByTransferTime(r.queue))
	n := float64(atomic.LoadInt64(&r.counter))

	if n > 1{
		r.bandwidth *= (n - 1) / n
	}
	if r.queue[0] == e {
		x := r.queue[0].remainingSize / r.bandwidth
		e.timeEnd = &x
	}
}

func (r *Resource) CounterAdd() {
	atomic.AddInt64(&r.counter, 1)
}

func (r *Resource) CounterMinus() {
	atomic.AddInt64(&r.counter, -1)
}
