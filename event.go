package main

import (
	"sync"
	"fmt"
)

type Resource struct {
	queue     []*Event
	bandwidth float64

	lastTimeRequest float64
	mutex           sync.Mutex
}

type Event struct {
	id        string
	size      float64
	timeStart float64
	timeEnd   float64

	remainingSize float64
	label         string
	callbacks     []func(*Event)
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