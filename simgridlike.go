package main

import "sync"

type RoutesMap map[Route]*Link

func (routeMap RoutesMap) Get(r Route) (d *Link, ok bool) {
	d, ok = routeMap[r]
	if ok {
		return
	}
	d, ok = routeMap[Route{start: r.finish, finish: r.start}]
	if ok {
		return
	}
	return nil, false
}

type Route struct {
	start  string
	finish string
}

type Link struct {
	*Resource
}

func(r *Resource) putEvents(events ...*Event){
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.queue = append(r.queue, events...)
}


func MSG_platform_init(env *Environment) {
	platform := make(map[Route]*Link)
	v := Route{start: "A",
		finish: "B"}
	platform[v] = &Link{
		bandwidth: 1,
		mutex : &sync.Mutex{},
	}
	env.routesMap = platform
}

func MSG_task_send(worker *Worker, sender string, address string, size float64) interface{} {
	wg := <- worker.resumeChan
	defer wg.Done()

	route := Route{sender, address}
	timeEnd := worker.env.currentTime + size / worker.env.routesMap[route].bandwidth

	eventA := Event{size: size,
		timeStart: worker.env.currentTime,
		timeEnd:   timeEnd,
		label:     sender}
	eventB := Event{size: size,
		timeStart: worker.env.currentTime,
		timeEnd:   timeEnd,
		label:     address}

	worker.env.PutEvents(&eventA, &eventB)

	// Should improve in the future!
	worker.env.routesMap[route].putEvents(&eventA, &eventB)
	return nil
}

func MSG_task_receive(worker *Worker) interface{} {
	wg := <- worker.resumeChan
	defer wg.Done()
	//_, value, _ := reflect.Select(env.cases)
	return nil
}

