package state

import "log"

func Init() {
	err := Load()
	if err != nil {
		log.Fatal("error loading printers", err)
	}

	go updateRoutine()
}

func updateRoutine() {
	for {
		Update()
	}
}
