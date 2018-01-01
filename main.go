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
	fmt.Println("end of sender 1")
	w.closeChan <-EndOfProcess{}
	fmt.Println("end of sender 2")
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
	fmt.Println("end of receiver 1")
	w.closeChan <- EndOfProcess{}
	fmt.Println("end of receiver 2")
}

func master(env *Environment, until float64, wg *sync.WaitGroup) {
	defer wg.Done()
	for env.currentTime < until {
		minTime := math.MaxFloat32
		var minChannel chan interface{}
		n := len(env.managerChannels)
		findEvent := false
		for i := 0; i < n; i++ {
			env.managerChannels[i].askChannel <- struct{}{}
			fmt.Println("heer")
			switch response := (<-env.managerChannels[i].answerChannel).(type) {
			case string:
			case float64:
				response = float64(response)
				if response < minTime {
					minTime = response
					minChannel = env.managerChannels[i].askChannel
				}
				findEvent = true
			}
			fmt.Println("heer")
		}
		if findEvent {
			minChannel <- true
			env.currentTime = minTime
			fmt.Println(minTime)
		}else {
			break
		}
	}
	n := len(env.managerChannels)
	for i := 0; i < n; i++ {
		close(env.closeChannels[i])
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
