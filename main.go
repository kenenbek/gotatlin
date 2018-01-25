package main

import (
	"fmt"
	"reflect"
	"sync"
)

func (w *Worker) send() {
	for i := float64(1); i < 3; i++ {
		w.MSG_task_send("GMAIL", 2*i)
	}
	<- w.resumeChan
	w.noMoreEvents = true
	w.resumeChan <- struct {}{}
}

func (w *Worker) receive() {
	for i := 1; i < 3; i++ {
		w.MSG_task_receive("GMAIL")
	}
	<- w.resumeChan
	w.noMoreEvents = true
	w.resumeChan <- struct {}{}
}

func master(env *Environment, until interface{}, wg *sync.WaitGroup) {
	if until != nil {
		switch until := until.(type) {
		case nil:
			//do nothing
		default:
			untilFloat64 := until.(float64)
			globalStop := ConstantEvent{
				Event: &Event{timeEnd: &untilFloat64},
			}
			globalStop.callbacks = append(globalStop.callbacks, env.stopSimulation, globalStop.print)
			env.PutEvents(&globalStop)
		}
	}
	// Initial
	var currentEvent EventInterface
	var isWorkerAlive bool
	defer wg.Done()

	cases := env.DoCases(env.workers...)
	env.SendCases(cases)
	env.WaitWorkers(cases)

	env.calculateTwinEvents()
	currentEvent, isWorkerAlive = env.Step()

	for !env.shouldStop {

		if isWorkerAlive{
			cases = env.DoCasesForEvent(currentEvent)
			env.SendCases(cases)
			env.WaitWorkers(cases)
		}

		env.calculateTwinEvents()
		currentEvent, isWorkerAlive = env.Step()
	}

	fmt.Println("end-master")
}

func main() {
	env := new(Environment)
	MSG_platform_init(env)
	link := make(chan float64)
	var wg sync.WaitGroup
	wg.Add(1)
	_ = NewWorkerSender(env, link, "A")
	_ = NewWorkerReceiver(env, link, "B")

	go master(env, float64(51.61), &wg)
	wg.Wait()
}

func (env *Environment) DoCases(workers ...*Worker) []reflect.SelectCase {
	cases := make([]reflect.SelectCase, len(workers))
	for i := range workers {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(workers[i].resumeChan)}
	}
	return cases
}

func (env *Environment) DoCasesForEvent(event EventInterface) []reflect.SelectCase {
	var cases []reflect.SelectCase
	switch event.(type) {
	case *TransferEvent:
		te := event.(*TransferEvent)
		cases = append(cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(te.worker.resumeChan)})
		cases = append(cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(te.twinEvent.worker.resumeChan)})
	case *ConstantEvent:
		ce := event.(*ConstantEvent)
		cases = append(cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ce.worker.resumeChan)})
	}
	return cases
}


func (env *Environment) SendCases(cases []reflect.SelectCase){
	for i := range cases {
		cases[i].Chan.Interface().(chan struct{}) <- struct{}{}
	}
}

func (env *Environment) WaitWorkers(cases []reflect.SelectCase) {
	remaining := len(cases)
	for remaining > 0 {
		reflect.Select(cases)
		remaining--
	}
}