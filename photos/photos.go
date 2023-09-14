// Copyright 2023 Matthew P. Dargan.
// SPDX-License-Identifier: Apache-2.0

package photos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// A Client is an HTTP client for interacting with the Photos API.
type Client struct {
	// Client is the HTTP client used to make requests to the Photos API.
	*http.Client

	// URL specifies the Photos endpoint. It defaults to https://jsonplaceholder.typicode.com/photos
	// but can be changed to localhost for testing purposes.
	URL string
}

const photosURL = "https://jsonplaceholder.typicode.com/photos"

// NewClient creates and returns a new Client with the given HTTP client.
func NewClient(client *http.Client) *Client {
	return &Client{Client: client, URL: photosURL}
}

// Photo represents details about a photo retrieved from the Photos API.
type Photo struct {
	AlbumID      int    `json:"albumId"`      // unique identifier of the album to which the photo belongs.
	ID           int    `json:"id"`           // unique identifier of the photo.
	Title        string `json:"title"`        // title of the photo.
	URL          string `json:"url"`          // URL to view the full-sized photo.
	ThumbnailURL string `json:"thumbnailUrl"` // URL to view the thumbnail version of the photo.
}

var (
	// ErrInvalidRequest is returned when the Photos API request is invalid.
	ErrInvalidRequest = errors.New("photos: invalid request")

	// ErrFailedRequest is returned when the Photos API request fails.
	ErrFailedRequest = errors.New("photos: failed to perform request")

	// ErrInvalidStatus is returned when the Photos API returns an invalid status code.
	ErrInvalidStatus = errors.New("photos: API request failed, status code")

	// ErrDecodeAPIResponse is returned when there is an error decoding the Photos API response body.
	ErrDecodeAPIResponse = errors.New("photos: failed to decode Photos API response body")
)

// Photos searches for photos in the Photos API and allows filtering by albumIds, ids, titles, urls, and/or thumbnailUrls.
//
// Example filters:
//   - Filtering by album ID: params["albumId"] = []string{"1"}
//   - Filtering by ID: params["id"] = []string{"1", "2"}
//   - Filtering by title: params["title"] = []string{"accusamus beatae ad facilis cum similique qui sunt"}
//   - Filtering by URL: params["url"] = []string{"https://via.placeholder.com/600/92c952"}
//   - Filtering by thumbnail URL: params["thumbnailUrl"] = []string{"https://via.placeholder.com/150/92c952"}
func (c *Client) Photos(ctx context.Context, params map[string][]string) ([]Photo, error) {
	req, err := newRequest(ctx, c.URL, params)
	if err != nil {
		return nil, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedRequest, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w %d", ErrInvalidStatus, resp.StatusCode)
	}
	var photos []Photo
	if err = json.NewDecoder(resp.Body).Decode(&photos); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDecodeAPIResponse, err)
	}
	return photos, nil
}

func newRequest(ctx context.Context, url string, params map[string][]string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}
	qry := req.URL.Query()
	for k, v := range params {
		for _, value := range v {
			qry.Add(k, value)
		}
	}
	req.URL.RawQuery = qry.Encode()
	return req, nil
}
