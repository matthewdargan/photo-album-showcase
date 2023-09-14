package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/matthewdargan/photo-album-showcase/photos"
)

const timeout = 5

var params = map[string][]string{
	"id":      {"1", "2", "3"},
	"albumId": {"1"},
}

func main() {
	pc := photos.NewClient(&http.Client{Timeout: time.Second * timeout})
	ps, err := pc.Photos(context.Background(), params)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(ps)
	// TODO: Create entry point into CLI application here
	// For example:
	// $ photos --id 1 2 3 --albumId 1
}
