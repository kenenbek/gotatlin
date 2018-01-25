package main

import (
	"sync"
)

type RoutesMap map[Route]*Link

func (routeMap RoutesMap) Get(r Route) (d *Link) {
	d, ok := routeMap[r]
	if ok {
		return
	}
	d, ok = routeMap[Route{start: r.finish, finish: r.start}]
	if ok {
		return
	}
	return nil
}

type Route struct {
	start  *Host
	finish *Host
}

type Link struct {
	*Resource
}

//func (r *Resource) putEvents(events ...*Event) {
//	r.mutex.Lock()
//	defer r.mutex.Unlock()
//	r.queue = append(r.queue, events...)
//}

func MSG_platform_init(env *Environment) {
	platform := make(map[Route]*Link)
	hostsMap := make(map[string]*Host)
	workersMap := make(map[string]*Worker)

	hostA := &Host{name: "A"}
	hostB := &Host{name: "B"}

	v := Route{start: hostA,
		finish: hostB,
	}

	platform[v] = &Link{
		Resource: &Resource{
			bandwidth: 1,
			mutex:     sync.Mutex{},
			queue:     []*TransferEvent{},
			counter:   0,
		},
	}
	hostsMap["A"] = hostA
	hostsMap["B"] = hostB

	env.routesMap = platform
	env.hostsMap = hostsMap
	env.workersMap = workersMap
}

func (worker *Worker) MSG_task_send(receiver string, size float64) interface{} {
	<-worker.resumeChan

	event := TransferEvent{
		Event: &Event{
			timeStart: worker.env.currentTime,
			worker:    worker,
			env:worker.env,
		},

		size:          size,
		remainingSize: size,

		sender:   worker.name,
		receiver: receiver,

		send: true,
		recv: false,

		}
	event.callbacks = append(event.callbacks, event.print)
	worker.env.PutEvents(&event)
	worker.resumeChan <- struct {}{}
	return nil
}

func (worker *Worker) MSG_task_receive(listener string) interface{} {
	<-worker.resumeChan

	event := TransferEvent{
		Event: &Event{
			worker: worker,
			env:worker.env,
		},
		listener: listener,
		send:     false,
		recv:     true,
	}
	event.callbacks = append(event.callbacks, event.print)
	worker.env.PutEvents(&event)
	worker.resumeChan <- struct {}{}
	return nil
}
