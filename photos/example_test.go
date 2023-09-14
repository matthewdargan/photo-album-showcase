// Copyright 2023 Matthew P. Dargan.
// SPDX-License-Identifier: Apache-2.0

package photos_test

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/matthewdargan/photo-album-showcase/photos"
)

func ExampleClient_Photos() {
	c := &http.Client{Timeout: time.Second * 5}
	client := photos.NewClient(c)
	params := map[string][]string{
		"id":      {"1", "2", "5"},
		"albumId": {"1"},
	}
	ps, err := client.Photos(context.Background(), params)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(ps)
	}
	// Output:
	// [{1 1 accusamus beatae ad facilis cum similique qui sunt https://via.placeholder.com/600/92c952 https://via.placeholder.com/150/92c952} {1 2 reprehenderit est deserunt velit ipsam https://via.placeholder.com/600/771796 https://via.placeholder.com/150/771796} {1 5 natus nisi omnis corporis facere molestiae rerum in https://via.placeholder.com/600/f66b97 https://via.placeholder.com/150/f66b97}]
}
