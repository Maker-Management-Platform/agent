package main

import (
	"log"

	"github.com/eduardooliveira/stLib/v2/web"
)

func main() {

	if err := web.RenderIndex("../../frontend/index.html"); err != nil {
		log.Fatalf("Error rendering index: %v", err)
	}
}
