package main

import (
	"fmt"
	"sync"
)

func ProcWrapper(env *Environment, processStrategy func(), w *Worker) {
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
					label:
					if len(w.queue) > 0 {
						w.answerChannel <- w.queue[0]
					} else if len(w.queue) == 0 && w.noMoreEvents {
						// No more events can happen and queue is empty
						w.answerChannel <- "E"
					} else if len(w.queue) == 0 && !w.noMoreEvents {
						for len(w.queue) == 0 {
							fmt.Println("start wait", w.noMoreEvents, w.name)
							w.cv.Wait()
							fmt.Println("end wait")
							if w.noMoreEvents == true {
								fmt.Println("wait and nomore events")
								goto label
							}
						}
						w.answerChannel <- w.queue[0]
					}
					w.cv.L.Unlock()
				}
			case <-w.closeChan:
				w.cv.L.Lock()
				w.noMoreEvents = true
				w.cv.L.Unlock()
			}
		}
		fmt.Println("end-worker1")
	}()
	w.queue = append(w.queue, 0)
	go processStrategy()
}

func NewWorkerReceiver(env *Environment, link chan float64) *Worker {
	w := &Worker{
		name:    "receiver",
		Process: NewProcess(env),
		env:     env,
		link:    link,
		cv:      sync.NewCond(&sync.Mutex{}),
		noMoreEvents:false}

	//w.queue = []float64{0.1, 0.3, 0.5}
	w.queue = []float64{}
	ProcWrapper(env, w.receive, w)
	return w
}

func NewWorkerSender(env *Environment, link chan float64) *Worker {
	w := &Worker{
		name:    "sender",
		Process: NewProcess(env),
		env:     env,
		link:    link,
		cv:      sync.NewCond(&sync.Mutex{}),
		noMoreEvents:false}
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
		askChannel:    ask,
		answerChannel: answer,
		closeChan:     closeChan,
	}
}
