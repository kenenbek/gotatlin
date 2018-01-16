package main

import (
	"sync"
	"sort"
)

type Environment struct {
	currentTime     float64
	managerChannels []pairChannel
	closeChannels   []chan interface{}
	workers         []*Worker
	routesMap       RoutesMap
	queue           []*Event
	mutex           sync.Mutex
	shouldStop      bool
}

func (env *Environment) stopSimulation(_ *Event) {
	env.shouldStop = true
}

func (env *Environment) PutEvents(events ...*Event) {
	env.mutex.Lock()
	defer env.mutex.Unlock()
	env.queue = append(env.queue, events...)
}

func (env *Environment) update(deltaTime float64) {
	for key := range env.routesMap {
		queue := env.routesMap[key].queue
		sort.Sort(ByTime(queue))
		bandwidthOld := env.routesMap[key].bandwidth / (float64(len(queue)) + 1)
		bandwidthNew := env.routesMap[key].bandwidth / (float64(len(queue)))
		for index := range queue{
			queue[index].remainingSize -= deltaTime * bandwidthOld
		}
		//update only for minimum event
		queue[0].timeEnd = env.currentTime + queue[0].remainingSize / bandwidthNew

	}
}
