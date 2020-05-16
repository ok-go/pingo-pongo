package main

import (
	"log"
	"net/http"
	"pingo_pongo"
)

func main() {
	server, err := pingo_pongo.NewPlayerServer()
	if err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatal(err)
	}
}
