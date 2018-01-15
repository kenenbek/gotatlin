package main

import (
	"fmt"
	"math"
	"sync"
)

func (w *Worker) send() {
	for i := float64(1); i < 10; i++ {
		MSG_task_send(w, "A", "B", 12*i)
	}
	w.cv.L.Lock()
	w.noMoreEvents = true
	w.cv.L.Unlock()
	w.cv.Signal()
}

func (w *Worker) receive() {
	x := Event{}
	for i := 1; i < 10; i++ {
		x = MSG_task_receive(w)
		w.cv.L.Lock()
		w.queue = append(w.queue, x)
		//fmt.Println("End receive", i)
		w.cv.L.Unlock()
		w.cv.Signal()
	}
	w.cv.L.Lock()
	w.noMoreEvents = true
	w.cv.L.Unlock()
	w.cv.Signal()
}

func master2(env *Environment, until float64, wg *sync.WaitGroup) {
	defer wg.Done()
	for env.currentTime < until {
		minTime := math.MaxFloat32
		n := len(env.workers)
		var indexOfProcess int
		findEvent := false
		for i := 0; i < n; i++ {
			env.workers[i].resumeChan <- struct {}{}
			
		}

		if findEvent {
			env.workers[indexOfProcess].deleteMinimumEvent()
			env.currentTime = minTime
			fmt.Println(minTime)
		} else {
			break
		}

		// Delete worker with noMore events
		j := 0
		for index := range env.workers {
			//_ := env.workers[index]
			if env.workers[index].hasMoreEvents() {
				env.workers[j] = env.workers[index]
				j++
			} else {
				//fmt.Println("end of", object.name, "object")
			}
		}
		env.workers = env.workers[:j]
	}

	fmt.Println("end-master")
}

func master(env *Environment, until float64, wg *sync.WaitGroup) {
	defer wg.Done()
	for env.currentTime < until {
		minTime := math.MaxFloat32
		n := len(env.workers)
		var indexOfProcess int
		findEvent := false
		for i := 0; i < n; i++ {
			timeEvent := env.workers[i].getMinimumEvent()
			switch timeEvent.(type) {
			case nil:
			case Event:
				timeEvent, _ := timeEvent.(Event)
				if timeEvent.timeEnd < minTime {
					minTime = timeEvent.timeEnd
					indexOfProcess = i
					findEvent = true
				}
			}

		}
		if findEvent {
			env.workers[indexOfProcess].deleteMinimumEvent()
			env.currentTime = minTime
			fmt.Println(minTime)
		} else {
			break
		}

		// Delete worker with noMore events
		j := 0
		for index := range env.workers {
			//_ := env.workers[index]
			if env.workers[index].hasMoreEvents() {
				env.workers[j] = env.workers[index]
				j++
			} else {
				//fmt.Println("end of", object.name, "object")
			}
		}
		env.workers = env.workers[:j]
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
	//_ = NewWorkerSender(env, link)
	_ = NewWorkerReceiver(env, link)

	go master(env, float64(51.61), &wg)
	wg.Wait()
}
