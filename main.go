package main

import (
	"fmt"
	"sync"
)

func (w *Worker) send() {
	for i := float64(1); i < 10; i++ {
		w.MSG_task_send("Recv", 12*i)
	}
	w.cv.L.Lock()
	w.noMoreEvents = true
	w.cv.L.Unlock()
}

func (w *Worker) receive() {
	for i := 1; i < 10; i++ {
		w.MSG_task_receive()
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
	var currentEvent *Event
	defer wg.Done()
	for !env.shouldStop {
		n := len(env.workers)

		var WaitGWorkers sync.WaitGroup
		WaitGWorkers.Add(len(env.workers))
		for i := 0; i < n; i++ {
			//fmt.Println("master")
			env.workers[i].resumeChan <- &WaitGWorkers
		}
		WaitGWorkers.Wait()

		// Pop from queue
		currentEvent = env.PopFromQueue()

		//Update duration
		env.updateQueue(currentEvent.timeEnd - env.currentTime)

		env.currentTime = currentEvent.timeEnd
		fmt.Println(env.currentTime)

		// Delete worker with noMore events
		j := 0
		for index := range env.workers {
			if env.workers[index].hasMoreEvents() {
				env.workers[j] = env.workers[index]
				j++
			}
		}
		env.workers = env.workers[:j]



		//
		for !env.shouldStop{
			currentEvent
		}
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
