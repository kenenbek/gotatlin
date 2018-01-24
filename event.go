package main

import (
	"fmt"
)

type EventInterface interface {
	update(float64)
	calculateTimeEnd()
	sendAble() bool
	receiveAble() bool
	getWorker() *Worker
	getTimeEnd() *float64
	getSize() float64
	getCallbacks() []func(EventInterface)
	print(EventInterface)
}

type Event struct {
	id        string
	timeStart float64
	timeEnd   *float64

	callbacks []func(EventInterface)
	worker    *Worker
	env       *Environment
}

type ConstantEvent struct {
	*Event
}

type TransferEvent struct {
	*Event
	size float64

	remainingSize float64
	resource      interface{}
	sender        string
	receiver      string
	listener      string

	recv bool
	send bool

	twinEvent *TransferEvent
}

func (e *ConstantEvent) update(_ float64) {
	//Nothing to update because event is immutable
}

func (e *ConstantEvent) sendAble() bool {
	return false
}

func (e *ConstantEvent) receiveAble() bool {
	return false
}

func (e *ConstantEvent) calculateTimeEnd() {

}

func (e *ConstantEvent) getWorker() *Worker {
	return e.worker
}

func (e *ConstantEvent) getTimeEnd() *float64 {
	return e.timeEnd
}

func (e *ConstantEvent) getSize() float64 {
	return 0
}

func (e *ConstantEvent) getCallbacks() []func(EventInterface) {
	return e.callbacks
}

func (e *ConstantEvent) print(_ EventInterface) {
	fmt.Printf("Start %v | End %v\n", e.timeStart, *e.timeEnd)
}

func (e *TransferEvent) update(deltaTime float64) {
	e.remainingSize -= (e.env.currentTime - deltaTime) * e.resource.(*Link).bandwidth
}

func (e *TransferEvent) calculateTimeEnd() {
	x := e.env.currentTime + e.remainingSize/(e.resource.(*Link).bandwidth/float64(e.resource.(*Link).counter))
	e.timeEnd = &x
}

func (e *TransferEvent) sendAble() bool {
	return e.send
}

func (e *TransferEvent) receiveAble() bool {
	return e.recv
}

func (e *TransferEvent) getListenAddress() string {
	return e.listener
}

func (e *TransferEvent) getReceiverAddress() string {
	return e.receiver
}

func (e *TransferEvent) getWorker() *Worker {
	return e.worker
}

func (e *TransferEvent) getTimeEnd() *float64 {
	return e.timeEnd
}

func (e *TransferEvent) getSize() float64 {
	return e.size
}

func (e *TransferEvent) getCallbacks() []func(EventInterface) {
	return e.callbacks
}

func (e *TransferEvent) print(_ EventInterface) {
	if e.timeEnd == nil{
		fmt.Printf("Start %v | End %v\n", e.timeStart)
	}else {
		fmt.Printf("Start %v | End %v\n", e.timeStart, *e.timeEnd)
	}
}

type ByTime []EventInterface

func (s ByTime) Len() int {
	return len(s)
}
func (s ByTime) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByTime) Less(i, j int) bool {
	if s[i].getTimeEnd() == nil && s[j].getTimeEnd() == nil {
		return s[i].getSize() < s[j].getSize()
	} else if s[i].getTimeEnd() == nil {
		return false
	} else if s[j].getTimeEnd() == nil {
		return true
	} else {
		return *(s[i].getTimeEnd()) < *(s[j].getTimeEnd())
	}
}

type ByTransferTime []*TransferEvent

func (s ByTransferTime) Len() int {
	return len(s)
}
func (s ByTransferTime) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByTransferTime) Less(i, j int) bool {
	if s[i].getTimeEnd() == nil && s[j].getTimeEnd() == nil {
		return s[i].getSize() < s[j].getSize()
	} else if s[i].getTimeEnd() == nil {
		return false
	} else if s[j].getTimeEnd() == nil {
		return true
	} else {
		return *(s[i].getTimeEnd()) < *(s[j].getTimeEnd())
	}
}

func (e *Event) String() string {
	return fmt.Sprintf("Start %v | End  %v\n", e.timeStart, e.timeEnd)
}
