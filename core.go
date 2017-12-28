package main

import (
	"fmt"
)

func ProcWrapper(env *Environment, procFn func(), w *Worker) {
	go func() {
		for {
			select {
			case rec := <-w.askChannel:
				switch rec.(type) {
				case bool:
					_, w.queue = w.queue[0], w.queue[1:]
				default:
					if len(w.queue) > 0 {
						w.answerChannel <- w.queue[0]
					} else {
						w.answerChannel <- "E"
					}
				}
			case <-w.closeChan:
				return
			}
		}
		fmt.Println("end-worker1")
	}()
	go procFn()
}


func NewWorkerReceiver(env *Environment, link chan float64) *Worker {
	w := &Worker{
		Process:NewProcess(env),
		env:env,
		link:link}

	w.queue = []float64{0.1, 0.3, 0.5}
	ProcWrapper(env, w.receive, w)
	return w
}


func NewWorkerSender(env *Environment, link chan float64) *Worker {
	w := &Worker{
		Process:NewProcess(env),
		env:env,
		link:link}
	w.queue = []float64{0.2, 0.4, 0.6}
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