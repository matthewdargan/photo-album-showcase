// Copyright 2023 Matthew P. Dargan.
// SPDX-License-Identifier: Apache-2.0

/*
Package photos provides a client for the [Photos API],
enabling request handling and response parsing.

To interact with the Photos API, create a [Client]:

	c := &http.Client{Timeout: time.Second * 5}
	client := photos.NewClient(c)
	params := map[string][]string{
		"id":      {"1", "2", "5"},
		"albumId": {"1"},
	}
	ps, err := client.Photos(context.Background(), params)
	if err != nil {
		// handle error
	}

For more details on the available methods and their usage,
see the examples under [Client].

[Photos API]: https://jsonplaceholder.typicode.com/photos
*/
package photos
