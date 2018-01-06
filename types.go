package main

import (
	"sync"
)

type Environment struct {
	currentTime     float64
	managerChannels []pairChannel
	closeChannels   []chan interface{}
	sliceOfProcesses []*Worker
}

type askChannel chan interface{}
type answerChannel chan interface{}

type pairChannel struct {
	askChannel
	answerChannel
}

type Process struct {
	env *Environment
	askChannel
	answerChannel
	closeChan chan interface{}
}

type Worker struct {
	*Process
	env *Environment
	link chan float64
	queue []float64
	name string
	cv *sync.Cond
	noMoreEvents bool
	mutex sync.RWMutex
}

type EndOfProcess struct {

}


func(w *Worker) getMinimumEvent() float64{
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	queueLength := len(w.queue)
	if queueLength == 0{
		return -42
	} else {
		return w.queue[0]
	}
}

func(w *Worker) deleteMinimumEvent(){
	w.mutex.Lock()
	defer w.mutex.Unlock()
	_, w.queue = w.queue[0], w.queue[1:]
}

func(w *Worker) hasMoreEvents() bool{
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	y := len(w.queue) > 0
	return y || (!w.noMoreEvents)
}