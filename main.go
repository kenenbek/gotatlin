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
	w.noMoreEvents = true
	fmt.Println("start signal of nomore events", w.name)
	w.cv.Signal()
	fmt.Println("end signal of nomore events", w.name)
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
	w.noMoreEvents = true
	fmt.Println("start signal of nomore events", w.name)
	w.cv.Signal()
	fmt.Println("end signal of nomore events", w.name)
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
			if timeEvent < minTime {
				minTime = timeEvent
				indexOfProcess = i
				findEvent = true
			}
		}
		if findEvent {
			fmt.Println("kotok", indexOfProcess, env.sliceOfProcesses[indexOfProcess].name)
			env.sliceOfProcesses[indexOfProcess].deleteMinimumEvent()
			env.currentTime = minTime
			fmt.Println(minTime)
		} else {
			break
		}

		// Delete worker with noMore events
		j := 0
		for index := range env.sliceOfProcesses {
			if env.sliceOfProcesses[index].hasMoreEvents() {
				env.sliceOfProcesses[j] = env.sliceOfProcesses[index]
				j++
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
