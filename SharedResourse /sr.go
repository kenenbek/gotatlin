package main

import (
	"fmt"
)

type SharedResource struct {
	queue []*Event
	bandwidth float64


	lastTimeRequest float64
}

type Event struct {
	id string
	size float64
	timeStart float64
	timeEnd float64

	remainingSize float64
}

func(resource *SharedResource) put(event *Event){
	if len(resource.queue) == 0 {
		event.timeEnd = event.size * float64(len(resource.queue)+1) / resource.bandwidth
		resource.queue = append(resource.queue, event)
		return
	}

	currentTime := event.timeStart
	for index:= range resource.queue{
		resource.queue[index].remainingSize
	}

}

func(resource *SharedResource) printEvents(){
	for index := range resource.queue{
		fmt.Println(resource.queue[index])
	}
}

func (event* Event) String() string {
	return fmt.Sprintf("Id: %v| Start: %v| End: %v", event.id, event.timeStart, event.timeEnd)
}
