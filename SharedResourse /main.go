package main

func main() {
	storage := SharedResource{queue: []*Event{}, bandwidth: 10}
	defer storage.printEvents()

	event1 := &Event{size: 10,
		timeStart: 0}
	event2 := &Event{size: 10,
		timeStart: 0}
	storage.put(event1)
	storage.put(event2)

}
