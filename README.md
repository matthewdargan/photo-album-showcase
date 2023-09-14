# Photo Album Showcase

[![GoDoc](https://godoc.org/github.com/matthewdargan/photo-album-showcase?status.svg)](https://godoc.org/github.com/matthewdargan/photo-album-showcase)
[![Build Status](https://github.com/matthewdargan/photo-album-showcase/actions/workflows/go-ci.yml/badge.svg?branch=main)](https://github.com/matthewdargan/photo-album-showcase/actions/workflows/go-ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/matthewdargan/photo-album-showcase)](https://goreportcard.com/report/github.com/matthewdargan/photo-album-showcase)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)

Photo-album-showcase fetches photos from the [Photos API](https://jsonplaceholder.typicode.com/photos). Results can be filtered by album IDs, IDs, titles, URLs, and/or thumbnail URLs.

Usage:

    photo-album-showcase [-id ids...] [-albumid albumids...] [-title 'titles...'] [-url urls...] [-thumburl thumburls...]

The `-id` flag filters photos by one or more photo IDs. Multiple IDs must be comma-separated.

The `-albumid` flag filters photos by one or more album IDs. Multiple album IDs must be comma-separated.

The `-title` flag filters photos by one or more titles. Multiple titles must be comma-separated. Titles containing spaces must be enclosed with single quotes (e.g., `-title 'title 1','title 2'`).

The `-url` flag filters photos by one or more URLs. Multiple URLs must be comma-separated.

The `-thumburl` flag filters photos by one or more thumbnail URLs. Multiple thumbnail URLs must be comma-separated.

## Installation

In order to install `photo-album-showcase`, install [Go](https://go.dev/doc/install) on your operating system. Once Go is installed, run the following to install `photo-album-showcase`:

```sh
go install github.com/matthewdargan/photo-album-showcase@latest
```

## Examples

```sh
# Fetch photos with IDs 1, 2, and 5 from album 1.
photo-album-showcase -id 1,2,5 -albumid 1

# Fetch photos with titles 'accusamus beatae ad facilis cum similique qui sunt' and 'reprehenderit est deserunt velit ipsam'.
photo-album-showcase -title 'accusamus beatae ad facilis cum similique qui sunt','reprehenderit est deserunt velit ipsam'

# Fetch photos with URLs https://via.placeholder.com/600/92c952 and https://via.placeholder.com/600/f66b97.
photo-album-showcase -url https://via.placeholder.com/600/92c952,https://via.placeholder.com/600/f66b97

# Fetch photos with thumbnail URLs https://via.placeholder.com/150/92c952 and https://via.placeholder.com/150/f66b97.
photo-album-showcase -thumburl https://via.placeholder.com/150/92c952,https://via.placeholder.com/150/f66b97
```
