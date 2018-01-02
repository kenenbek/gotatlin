package main

import "sync"

type Environment struct {
	currentTime     float64
	managerChannels []pairChannel
	closeChannels   []chan interface{}
}

type askChannel chan interface{}
type answerChannel chan interface{}

type pairChannel struct {
	askChannel
	answerChannel
}

type Process struct {
	env *Environment
	askChannel
	answerChannel
	closeChan chan interface{}
}

type Worker struct {
	*Process
	env *Environment
	link chan float64
	queue []float64
	name string
	cv *sync.Cond
}

type EndOfProcess struct {

}
