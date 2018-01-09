package main

import (
	"fmt"
	"math"
	"sync"
)

func (w *Worker) send() {
	for i := float64(1); i < 3; i++ {
		w.link <- i
	}
	w.cv.L.Lock()
	w.noMoreEvents = true
	w.cv.L.Unlock()
	w.cv.Signal()
}

func (w *Worker) receive() {
	x := float64(0)
	for i := 1; i < 3; i++ {
		x = <-w.link
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

func master(env *Environment, until float64, wg *sync.WaitGroup) {
	defer wg.Done()
	for env.currentTime < until {
		minTime := math.MaxFloat32
		n := len(env.sliceOfProcesses)
		var indexOfProcess int
		findEvent := false
		for i := 0; i < n; i++ {
			timeEvent := env.sliceOfProcesses[i].getMinimumEvent()
			switch timeEvent.(type) {
			case nil:
			case float64:
				timeEvent, _ := timeEvent.(float64)
				if timeEvent < minTime {
					minTime = timeEvent
					indexOfProcess = i
					findEvent = true
				}
			}
		}
		if findEvent {
			env.sliceOfProcesses[indexOfProcess].deleteMinimumEvent()
			env.currentTime = minTime
			fmt.Println(minTime)
		} else {
			break
		}

		// Delete worker with noMore events
		j := 0
		for index := range env.sliceOfProcesses {
			object := env.sliceOfProcesses[index]
			if env.sliceOfProcesses[index].hasMoreEvents() {
				env.sliceOfProcesses[j] = env.sliceOfProcesses[index]
				j++
			}else{
				fmt.Println("end of", object.name, "object")
			}
		}
		env.sliceOfProcesses = env.sliceOfProcesses[:j]
	}

	fmt.Println("end-master")
}

func main() {
	link := make(chan float64)
	var wg sync.WaitGroup
	wg.Add(1)
	env := new(Environment)
	_ = NewWorkerSender(env, link)
	_ = NewWorkerReceiver(env, link)

	go master(env, float64(51.61), &wg)
	wg.Wait()
}
