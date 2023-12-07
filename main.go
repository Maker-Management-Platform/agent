package main

import (
	"log"

	stlib "github.com/eduardooliveira/stLib/core"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	stlib.Run()
}
