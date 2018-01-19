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
	// Some amount of data has been sent over time
	for index := range env.queue {
		event := env.queue[index]
		event.remainingSize -= deltaTime * event.resource.(*Resource).bandwidth/float64(event.resource.(*Link).counter)
	}
	firstEvent := env.queue[0]

	switch firstEvent.resource.(type) {
	case Link, *Link:
		firstEvent.timeEnd = env.currentTime + firstEvent.remainingSize/(firstEvent.resource.(*Link).bandwidth/float64(firstEvent.resource.(*Link).counter))
	}
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
			event := env.queue[index]
			ReceiverSendersMap[event] = []*Event{}
			EventByNameMap[event.listener] = event
		}
	}
	for index := range env.queue {
		event := env.queue[index]
		if event.send {
			ReceiverSendersMap[EventByNameMap[event.receiver]] = append(ReceiverSendersMap[EventByNameMap[event.receiver]], event)
		}
	}

	for receiveEvent := range ReceiverSendersMap {
		for index := range ReceiverSendersMap[receiveEvent] {
			sendEvent := ReceiverSendersMap[receiveEvent][index]
			route := Route{receiveEvent.worker.host, sendEvent.worker.host}
			resource := env.routesMap.Get(route)
			sendEvent.resource = resource
			resource.Put(sendEvent)
		}
		sort.Sort(ByTime(ReceiverSendersMap[receiveEvent]))
		receiveEvent.twinEvent, ReceiverSendersMap[receiveEvent][0].twinEvent = ReceiverSendersMap[receiveEvent][0].twinEvent, receiveEvent.twinEvent
	}
	return nil
}


func (env *Environment) PopFromQueue() *Event {
	var currentEvent *Event
	//Sorting of events
	sort.Sort(ByTime(env.queue))

	currentEvent, env.queue = env.queue[0], env.queue[1:]

	if currentEvent.twinEvent != nil{
		for i := range env.queue{
			if currentEvent.twinEvent == env.queue[i]{
				copy(env.queue[i:], env.queue[i+1:])
				env.queue[len(env.queue)-1] = nil
				env.queue = env.queue[:len(env.queue)-1]

				// Process the event callbacks
				callbacks := currentEvent.twinEvent.callbacks
				currentEvent.twinEvent.callbacks = nil
				for _, callback := range callbacks {
					callback(currentEvent.twinEvent)
				}
				break
			}
		}
	}

	// Process the event callbacks
	callbacks := currentEvent.callbacks
	currentEvent.callbacks = nil
	for _, callback := range callbacks {
		callback(currentEvent)
	}
	return currentEvent
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
