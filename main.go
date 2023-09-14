// Copyright 2023 Matthew P. Dargan.
// SPDX-License-Identifier: Apache-2.0

// Photo-album-showcase fetches photos from the Photos API (https://jsonplaceholder.typicode.com/photos).
// Results can be filtered by album IDs, IDs, titles, URLs, and/or thumbnail URLs.
//
// Usage:
//
//	photo-album-showcase [-id ids...] [-albumid albumids...] [-title 'titles...'] [-url urls...] [-thumburl thumburls...]
//
// The -id flag filters photos by one or more photo IDs. Multiple IDs must be comma-separated.
//
// The -albumid flag filters photos by one or more album IDs. Multiple album IDs must be comma-separated.
//
// The -title flag filters photos by one or more titles. Multiple titles must be comma-separated.
// Titles containing spaces must be enclosed with single quotes (e.g., `-title 'title 1','title 2'`).
//
// The -url flag filters photos by one or more URLs. Multiple URLs must be comma-separated.
//
// The -thumburl flag filters photos by one or more thumbnail URLs. Multiple thumbnail URLs must be comma-separated.
//
// Examples:
//
// Fetch photos with IDs 1, 2, and 5 from album 1:
//
//	$ photo-album-showcase -id 1,2,5 -albumid 1
//
// Fetch photos with titles 'accusamus beatae ad facilis cum similique qui sunt' and 'reprehenderit est deserunt velit ipsam':
//
//	$ photo-album-showcase -title 'accusamus beatae ad facilis cum similique qui sunt','reprehenderit est deserunt velit ipsam'
//
// Fetch photos with URLs https://via.placeholder.com/600/92c952 and https://via.placeholder.com/600/f66b97:
//
//	$ photo-album-showcase -url https://via.placeholder.com/600/92c952,https://via.placeholder.com/600/f66b97
//
// Fetch photos with thumbnail URLs https://via.placeholder.com/150/92c952 and https://via.placeholder.com/150/f66b97:
//
//	$ photo-album-showcase -thumburl https://via.placeholder.com/150/92c952,https://via.placeholder.com/150/f66b97
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/matthewdargan/photo-album-showcase/photos"
)

const (
	stderr  = 2
	timeout = 5
)

var (
	id       = flag.String("id", "", "filter photos by ID(s), comma-separated")
	albumID  = flag.String("albumid", "", "filter photos by album ID(s), comma-separated")
	title    = flag.String("title", "", "filter photos by title(s), comma-separated, use single quotes for titles with spaces (e.g., -title 'title 1','title 2')")
	url      = flag.String("url", "", "filter photos by URL(s), comma-separated")
	thumbURL = flag.String("thumburl", "", "filter photos by thumbnail URL(s), comma-separated")
	help     = flag.Bool("h", false, "display usage")
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: photo-album-showcase [-id ids...] [-albumid albumids...] [-title 'titles...'] [-url urls...] [-thumburl thumburls...]\n")
	flag.PrintDefaults()
	os.Exit(stderr)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if *help {
		flag.Usage()
		return
	}

	params := make(map[string][]string)
	if *id != "" {
		params["id"] = strings.Split(*id, ",")
	}
	if *albumID != "" {
		params["albumId"] = strings.Split(*albumID, ",")
	}
	if *title != "" {
		params["title"] = strings.Split(*title, ",")
	}
	if *url != "" {
		params["url"] = strings.Split(*url, ",")
	}
	if *thumbURL != "" {
		params["thumbnailUrl"] = strings.Split(*thumbURL, ",")
	}

	client := photos.NewClient(&http.Client{Timeout: time.Second * timeout})
	ps, err := client.Photos(context.Background(), params)
	if err != nil {
		log.Fatal(err)
	}
	jsonData, err := json.MarshalIndent(ps, "", "  ")
	if err != nil {
		log.Fatalf("Failed to encode photos as JSON: %v", err)
	}
	fmt.Println(string(jsonData))
}
