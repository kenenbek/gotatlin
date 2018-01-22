package main

import (
	"sync"
)

func ProcWrapper(env *Environment, processStrategy func(), w *Worker) {
	//event := Event{timeStart: 0,
	//	timeEnd: 0}
	//
	//env.PutEvents(&event)
	go processStrategy()
}

func NewWorkerReceiver(env *Environment, link chan float64, hostName string) *Worker {
	host := env.getHostByName(hostName)
	w := &Worker{
		name:         "receiver",
		Process:      NewProcess(env, host),
		env:          env,
		link:         link,
		cv:           sync.NewCond(&sync.Mutex{}),
		noMoreEvents: false,
		mutex:        sync.RWMutex{}}

	ProcWrapper(env, w.receive, w)
	env.workers = append(env.workers, w)
	return w
}

func NewWorkerSender(env *Environment, link chan float64, hostName string) *Worker {
	host := env.getHostByName(hostName)
	w := &Worker{
		name:         "sender",
		Process:      NewProcess(env, host),
		env:          env,
		link:         link,
		cv:           sync.NewCond(&sync.Mutex{}),
		noMoreEvents: false,
		mutex:        sync.RWMutex{}}
	//w.queue = []float64{0.2, 0.4, 0.6}
	w.queue = []Event{}
	ProcWrapper(env, w.send, w)
	env.workers = append(env.workers, w)
	return w
}

func NewProcess(env *Environment, host *Host) *Process {
	ask := make(chan interface{})
	answer := make(chan interface{})
	closeChan := make(chan interface{})
	noMEC := make(chan bool)

	pairChan := pairChannel{ask, answer}
	env.managerChannels = append(env.managerChannels, pairChan)
	env.closeChannels = append(env.closeChannels, closeChan)
	return &Process{
		askChannel:       ask,
		answerChannel:    answer,
		waitEventsOrDone: closeChan,
		noMoreEventsChan: noMEC,
		resumeChan:       make(chan *sync.WaitGroup),
		host:             host,
	}
}
