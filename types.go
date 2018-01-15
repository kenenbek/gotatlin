package main

import (
	"reflect"
	"sync"
)

type Environment struct {
	currentTime      float64
	managerChannels  []pairChannel
	closeChannels    []chan interface{}
	workers []*Worker
	platform         map[Route]Link
	cases            []reflect.SelectCase
	queue            []Event
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
	waitEventsOrDone chan interface{}
	noMoreEventsChan chan bool
}

type Worker struct {
	*Process
	env          *Environment
	link         chan float64
	queue        []Event
	name         string
	cv           *sync.Cond
	noMoreEvents bool
	mutex        sync.RWMutex
	resumeChan   chan interface{}
}

type EndOfProcess struct {
}

func (w *Worker) getMinimumEvent() interface{} {
	w.cv.L.Lock()
	defer w.cv.L.Unlock()
	queueLength := len(w.queue)
	if queueLength == 0 && w.noMoreEvents {
		return nil
	} else if queueLength == 0 && !w.noMoreEvents {
		for len(w.queue) == 0 {
			w.cv.Wait()
			if len(w.queue) != 0 {
				return w.queue[0]
			}
			if w.noMoreEvents == true {
				return nil
			}
		}
		return w.queue[0]
	} else {
		return w.queue[0]
	}
	return nil
}

func (w *Worker) deleteMinimumEvent() {
	w.cv.L.Lock()
	defer w.cv.L.Unlock()
	_, w.queue = w.queue[0], w.queue[1:]
}

func (w *Worker) hasMoreEvents() bool {
	w.cv.L.Lock()
	defer w.cv.L.Unlock()
	y := len(w.queue) > 0
	return y || (!w.noMoreEvents)
}


func (worker *Worker) doWork(){
	<-worker.resumeChan
	worker.resumeChan <- struct {}{}
}