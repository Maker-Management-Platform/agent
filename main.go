package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	stlib "github.com/eduardooliveira/stLib/core"
)

func main() {
	go func() {
		http.ListenAndServe("localhost:8080", nil)
	}()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	stlib.Run()
}
