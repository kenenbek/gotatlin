package main

import (
	//"sort"
	"sync"
	//"fmt"
	"fmt"
	"sort"
)

type Host struct {
	name    string
	workers []*Worker
}

type Environment struct {
	currentTime     float64
	managerChannels []pairChannel
	closeChannels   []chan interface{}
	workers         []*Worker
	routesMap       RoutesMap
	queue           []*Event
	mutex           sync.Mutex
	shouldStop      bool
	hostsMap        map[string]*Host
	workersMap      map[string]*Worker
}

func (env *Environment) stopSimulation(_ *Event) {
	env.shouldStop = true
}

func (env *Environment) PutEvents(events ...*Event) {
	env.mutex.Lock()
	defer env.mutex.Unlock()
	env.queue = append(env.queue, events...)
}

func (env *Environment) updateQueue(deltaTime float64) {
	for index := range env.queue {
		env.queue[index].update(deltaTime)
	}
	firstElement := env.queue[0]

	switch firstElement.resource.(type) {
	case Link, *Link:
		firstElement.timeEnd = env.currentTime + firstElement.remainingSize/(firstElement.resource.(*Link).bandwidth/float64(firstElement.resource.(*Link).counter))
	}
}

func (env *Environment) PopFromQueue() *Event {
	var currentEvent *Event
	//Sorting of events
	sort.Sort(ByTime(env.queue))

	currentEvent, env.queue = env.queue[0], env.queue[1:]

	// Process the event callbacks
	callbacks := currentEvent.callbacks
	currentEvent.callbacks = nil
	for _, callback := range callbacks {
		callback(currentEvent)
	}
	return currentEvent
}

func (env *Environment) getHostByName(name string) *Host {
	return env.hostsMap[name]
}

func (env *Environment) getWorkerByName(name string) *Worker {
	return env.workersMap[name]
}

func (env *Environment) calculateTwinEvents(name string) interface{} {
	/*
		nonoptimal solution
	*/
	// receiver -- senders map
	ReceiverSendersMap := make(map[*Event][]*Event)
	EventByNameMap := make(map[string]*Event)
	for index := range env.queue {
		if env.queue[index].recv {
			ReceiverSendersMap[env.queue[index]] = []*Event{}
			EventByNameMap[env.queue[index].listener] = env.queue[index]
		}
	}
	for index := range env.queue {
		if env.queue[index].send {
			ReceiverSendersMap[EventByNameMap[env.queue[index].receiver]] = append(ReceiverSendersMap[EventByNameMap[env.queue[index].receiver]], env.queue[index])
		}
	}

	for receiveEvent := range ReceiverSendersMap {
		for index := range ReceiverSendersMap[receiveEvent] {
			route := Route{receiveEvent.worker.host, ReceiverSendersMap[receiveEvent][index].worker.host}
			env.routesMap.Get(route).Put(ReceiverSendersMap[receiveEvent][index])
		}
		sort.Sort(ByTime(ReceiverSendersMap[receiveEvent]))
		receiveEvent.twinEvent, ReceiverSendersMap[receiveEvent][0].twinEvent = ReceiverSendersMap[receiveEvent][0].twinEvent, receiveEvent.twinEvent
	}
	return nil
}

func (env *Environment) Step() {
	currentEvent := env.PopFromQueue()

	//Update duration
	env.updateQueue(currentEvent.timeEnd.(float64) - env.currentTime)

	env.currentTime = currentEvent.timeEnd.(float64)
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
}
