package main

import (
	"log"

	"github.com/andrewclowes/fresh8/importer/event"
)

func main() {
	p := event.NewPipeline()
	in := p.Run(nil)
	for r := range in {
		log.Println(r)
	}
	log.Println("End")
}
