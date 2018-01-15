package main

import (
	"reflect"
)

type Route struct {
	start  string
	finish string
}

type Link struct {
	linkChan  chan Event
	bandwidth float64

	resource *Resource
}

func MSG_platform_init(env *Environment) {
	platform := make(map[Route]Link)
	v := Route{start: "A",
		finish: "B"}
	platform[v] = Link{
		linkChan:  make(chan Event),
		bandwidth: 1,
	}
	env.platform = platform

	/*for reflection. receive from multiple goroutines*/
	cases := make([]reflect.SelectCase, 1)
	cases[0] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(env.platform[v].linkChan)}
	env.cases = cases
}

func MSG_task_send(worker *Worker, sender string, address string, size float64) interface{} {
	<- worker.resumeChan
	route := Route{sender, address}
	timeEnd := worker.env.currentTime + size/worker.env.platform[route].bandwidth

	eventA := Event{size: size,
		timeStart: worker.env.currentTime,
		timeEnd:   timeEnd,
		label:     sender}
	eventB := Event{size: size,
		timeStart: worker.env.currentTime,
		timeEnd:   timeEnd,
		label:     address}

	//worker.env.platform[route].linkChan <- event
	worker.env.queue = append(worker.env.queue, eventA, eventB)
	return nil
}

func MSG_task_receive(worker *Worker, env *Environment) interface{} {
	<- worker.resumeChan
	//_, value, _ := reflect.Select(env.cases)
	return nil
}

