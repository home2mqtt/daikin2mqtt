package main

import (
	"log"
	"net/url"

	"github.com/samthor/daikin-go/api"
)

func main() {
	devices := make(map[string]url.Values)

	disc, err := api.NewDiscover()
	if err != nil {
		log.Fatal(err)
	}

	err = disc.Announce()
	if err != nil {
		log.Fatal(err)
	}

	defer disc.Close()

	for id, v, err := disc.Next(); id != ""; {
		if err != nil {
			log.Println(err)
		} else {
			_, exists := devices[id]
			if !exists {
				log.Printf("%s (%v)", id, v)
				devices[id] = v
			}
			disc.Announce()
		}
	}
}
