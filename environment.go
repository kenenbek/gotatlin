package main

import (
	"sync"
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
	queue           []EventInterface
	mutex           sync.Mutex
	shouldStop      bool
	hostsMap        map[string]*Host
	workersMap      map[string]*Worker
}

func (env *Environment) stopSimulation(_ EventInterface) {
	env.shouldStop = true
}

func (env *Environment) PutEvents(events ...EventInterface) {
	env.mutex.Lock()
	defer env.mutex.Unlock()
	env.queue = append(env.queue, events...)
}

func (env *Environment) updateQueue(deltaTime float64) {
	// Some amount of data has been sent over time
	for index := range env.queue {
		event := env.queue[index]
		event.update(deltaTime)
	}
	firstEvent := env.queue[0]

	switch firstEvent.(type) {
	case *TransferEvent:
		firstEvent.calculateTimeEnd()
	}
}

func (env *Environment) getHostByName(name string) *Host {
	return env.hostsMap[name]
}

func (env *Environment) getWorkerByName(name string) *Worker {
	return env.workersMap[name]
}

func (env *Environment) calculateTwinEvents() interface{} {
	/*
		nonoptimal solution
	*/
	// receiver -- senders map
	ReceiverSendersMap := make(map[EventInterface][]EventInterface)
	EventByNameMap := make(map[string]EventInterface)

	for index := range env.queue {
		if env.queue[index].receiveAble() {
			event := env.queue[index].(*TransferEvent)
			ReceiverSendersMap[event] = []EventInterface{}
			EventByNameMap[event.getListenAddress()] = event
		}
	}

	for index := range env.queue {
		switch env.queue[index].(type) {
		case *TransferEvent:
			event := env.queue[index].(*TransferEvent)
			if event.sendAble() {
				ReceiverSendersMap[EventByNameMap[event.getReceiverAddress()]] = append(ReceiverSendersMap[EventByNameMap[event.getReceiverAddress()]], event)
			}
		}
	}

	for recEvent := range ReceiverSendersMap {
		receiveEvent := recEvent.(*TransferEvent)
		for index := range ReceiverSendersMap[receiveEvent] {
			sendEvent := ReceiverSendersMap[receiveEvent][index].(*TransferEvent)
			route := Route{receiveEvent.worker.host, sendEvent.worker.host}
			resource := env.routesMap.Get(route)
			sendEvent.resource = resource
			resource.Put(sendEvent)

		}
		sort.Sort(ByTime(ReceiverSendersMap[receiveEvent]))
		receiveEvent.twinEvent, ReceiverSendersMap[receiveEvent][0].(*TransferEvent).twinEvent = ReceiverSendersMap[receiveEvent][0].(*TransferEvent), receiveEvent
		receiveEvent.resource = receiveEvent.twinEvent.resource
		receiveEvent.timeEnd = receiveEvent.twinEvent.timeEnd
	}
	return nil
}

func (env *Environment) PopFromQueue() EventInterface {
	var currentEvent EventInterface
	//Sorting of events
	sort.Sort(ByTime(env.queue))

	currentEvent, env.queue = env.queue[0], env.queue[1:]
	switch currentEvent.(type) {
	case *TransferEvent:
		CE := currentEvent.(*TransferEvent)
		for i := range env.queue {
			if CE.twinEvent == env.queue[i] {
				copy(env.queue[i:], env.queue[i+1:])
				env.queue[len(env.queue)-1] = nil
				env.queue = env.queue[:len(env.queue)-1]

				// Process the event of twin event callbacks
				callbacks := CE.twinEvent.callbacks
				CE.twinEvent.callbacks = nil
				for _, callback := range callbacks {
					callback(CE.twinEvent)
				}
				break
			}
		}
	case *ConstantEvent:
	}

	// Process the event callbacks
	callbacks := currentEvent.getCallbacks()
	for _, callback := range callbacks {
		callback(currentEvent)
	}
	return currentEvent
}

func (env *Environment) Step() (EventInterface, bool) {
	currentEvent := env.PopFromQueue()

	if len(env.queue) > 0 {
		env.updateQueue(*currentEvent.getTimeEnd() - env.currentTime)

		// Delete worker with noMore events
		j := 0
		for index := range env.workers {
			if env.workers[index].hasMoreEvents() {
				env.workers[j] = env.workers[index]
				j++
			} else {
				env.workers[index] = nil
			}
		}
		env.workers = env.workers[:j]
	}else {
		env.shouldStop = true
	}
	env.currentTime = *currentEvent.getTimeEnd()
	return currentEvent, currentEvent.getWorker() != nil
}
