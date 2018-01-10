package main

import (
	"reflect"
)

type Route struct {
	start  string
	finish string
}

type Link struct {
	linkChan  chan float64
	bandwidth float64
}

func MSG_platform_init(env *Environment) {
	platform := make(map[Route]Link)
	v := Route{start: "A",
		finish: "B"}
	platform[v] = Link{
		linkChan:  make(chan float64),
		bandwidth: 1,
	}
	env.platform = platform

	/*for reflection. receive from multiple goroutines*/
	cases := make([]reflect.SelectCase, 1)
	cases[0] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(env.platform[v].linkChan)}
	env.cases = cases
}

func MSG_task_send(env *Environment, sender string, address string, size float64) {
	route := Route{sender, address}
	env.platform[route].linkChan <- size
	//bandwidth := env.platform[route]
}

func MSG_task_receive(env *Environment) float64{
	_, value, _ := reflect.Select(env.cases)
	return value.Float()
}
