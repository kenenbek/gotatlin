package main

import (
	"fmt"
	"sync"
)

func (w *Worker) send() {
	for i := float64(1); i < 10; i++ {
		w.MSG_task_send("GMAIL", 12*i)
	}
	w.cv.L.Lock()
	w.noMoreEvents = true
	w.cv.L.Unlock()
}

func (w *Worker) receive() {
	for i := 1; i < 10; i++ {
		w.MSG_task_receive("GMAIL")
	}
	w.cv.L.Lock()
	w.noMoreEvents = true
	w.cv.L.Unlock()
}

func master(env *Environment, until interface{}, wg *sync.WaitGroup) {
	if until != nil {
		switch until := until.(type) {
		default:
			untilFloat64 := until.(float64)
			globalStop := Event{timeEnd: untilFloat64}
			globalStop.callbacks = append(globalStop.callbacks, env.stopSimulation)
			env.PutEvents(&globalStop)
		}
	}
	// Initial
	var currentEvent *Event
	defer wg.Done()

	n := len(env.workers)

	var WaitGWorkers sync.WaitGroup
	WaitGWorkers.Add(len(env.workers))
	for i := 0; i < n; i++ {
		env.workers[i].resumeChan <- &WaitGWorkers
	}
	WaitGWorkers.Wait()

	env.Step()

	for !env.shouldStop {
		var singleWG sync.WaitGroup
		singleWG.Add(1)
		currentEvent.worker.resumeChan <- &singleWG
		singleWG.Wait()

		env.Step()
	}

	fmt.Println("end-master")
}

func main() {
	env := new(Environment)
	MSG_platform_init(env)
	link := make(chan float64)
	var wg sync.WaitGroup
	wg.Add(1)
	_ = NewWorkerSender(env, link)
	_ = NewWorkerReceiver(env, link)

	go master(env, float64(51.61), &wg)
	wg.Wait()
}
