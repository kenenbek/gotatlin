package main

import (
	"fmt"
	"sync"
)

func ProcWrapper(env *Environment, processStrategy func(), w *Worker) {
	noMoreEvents := false
	go func() {
		for {
			select {
			case rec := <-w.askChannel:
				switch rec.(type) {
				case bool:
					w.cv.L.Lock()
					_, w.queue = w.queue[0], w.queue[1:]
					w.cv.L.Unlock()
				default:
					w.cv.L.Lock()
					if len(w.queue) > 0 {
						w.answerChannel <- w.queue[0]
					} else if len(w.queue) == 0 && noMoreEvents{
						// No more events can happen and queue is empty
						w.answerChannel <- "E"
						w.cv.L.Unlock()
						fmt.Println("No more events can happen and queue is empty")
						return
					} else if len(w.queue) == 0 && !noMoreEvents{
						for len(w.queue) == 0  {
							w.cv.Wait()
						}
						w.answerChannel <- w.queue[0]
					}
					w.cv.L.Unlock()
				}
			case <-w.closeChan:
				w.cv.L.Lock()
				noMoreEvents = true
				fmt.Println("No more events")
				w.cv.L.Unlock()
			}
		}
		fmt.Println("end-worker1")
	}()
	go processStrategy()
}


func NewWorkerReceiver(env *Environment, link chan float64) *Worker {
	w := &Worker{
		Process:NewProcess(env),
		env:env,
		link:link,
		cv:sync.NewCond(&sync.Mutex{})}

	//w.queue = []float64{0.1, 0.3, 0.5}
	w.queue = []float64{}
	ProcWrapper(env, w.receive, w)
	return w
}


func NewWorkerSender(env *Environment, link chan float64) *Worker {
	w := &Worker{
		Process:NewProcess(env),
		env:env,
		link:link,
		cv:sync.NewCond(&sync.Mutex{})}
	//w.queue = []float64{0.2, 0.4, 0.6}
	w.queue = []float64{}
	ProcWrapper(env, w.send, w)
	return w
}

func NewProcess(env *Environment) *Process {
	ask := make(chan interface{})
	answer := make(chan interface{})
	closeChan := make(chan interface{})

	pairChan := pairChannel{ask, answer}
	env.managerChannels = append(env.managerChannels, pairChan)
	env.closeChannels = append(env.closeChannels, closeChan)
	return &Process{
		askChannel: ask,
		answerChannel: answer,
		closeChan: closeChan,
	}
}