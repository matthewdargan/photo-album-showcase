// Copyright 2023 Matthew P. Dargan.
// SPDX-License-Identifier: Apache-2.0

package photos_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/matthewdargan/photo-album-showcase/photos"
)

func TestPhotos(t *testing.T) {
	t.Parallel()
	t.Run("can find all photos", func(t *testing.T) {
		t.Parallel()
		f, err := os.Open("photos.json")
		if err != nil {
			t.Fatalf("Failed to open test data file: %v", err)
		}
		defer f.Close()
		var testData []photos.Photo
		if err := json.NewDecoder(f).Decode(&testData); err != nil {
			t.Fatalf("Failed to decode test data: %v", err)
		}
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(testData); err != nil {
				t.Fatalf("Failed to encode test data: %v", err)
			}
		}))
		defer ts.Close()

		// Create a new client with the mock server's URL.
		client := photos.NewClient(ts.Client())
		client.URL = ts.URL

		ps, err := client.Photos(context.Background(), nil)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(ps) != len(testData) {
			t.Errorf("Expected %d photos, got %d", len(testData), len(ps))
		}
	})

	t.Run("returns error if the client request was not successful", func(t *testing.T) {
		t.Parallel()
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()

		client := photos.NewClient(ts.Client())
		client.URL = ts.URL
		_, err := client.Photos(context.Background(), nil)
		want := fmt.Errorf("%w %d", photos.ErrInvalidStatus, http.StatusInternalServerError)
		if err == nil {
			t.Errorf("Expected an error, got nil")
		} else if err.Error() != want.Error() {
			t.Errorf("got %v, want %v", err, want)
		}
	})

	t.Run("returns error if the client cannot establish a connection", func(t *testing.T) {
		t.Parallel()
		client := photos.NewClient(http.DefaultClient)
		client.URL = "http://localhost"
		_, err := client.Photos(context.Background(), nil)
		want := fmt.Errorf("%w: %s", photos.ErrFailedRequest, `Get "http://localhost": dial tcp`)
		if err == nil {
			t.Errorf("Expected an error, got nil")
		} else if !strings.HasPrefix(err.Error(), want.Error()) {
			t.Errorf("expected error with prefix %q, got: %q", want.Error(), err.Error())
		}
	})

	t.Run("returns error if the response cannot be parsed into photos", func(t *testing.T) {
		t.Parallel()
		badData := `[123.1, 234.2]`
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(badData)); err != nil {
				t.Fatalf("Failed to encode test data: %v", err)
			}
		}))
		defer ts.Close()

		client := photos.NewClient(ts.Client())
		client.URL = ts.URL
		_, err := client.Photos(context.Background(), nil)
		want := fmt.Errorf("%w: %s", photos.ErrDecodeAPIResponse, "json: cannot unmarshal number into Go value of type photos.Photo")
		if err == nil {
			t.Errorf("Expected an error, got nil")
		} else if err.Error() != want.Error() {
			t.Errorf("got %v, want %v", err, want)
		}
	})
}

func TestPhotosIntegration(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("skipping integration tests")
	}
	testCases := []struct {
		name          string
		params        map[string][]string
		expectedCount int
	}{
		{
			name:          "can find photos with empty params",
			params:        map[string][]string{},
			expectedCount: 5000,
		},
		{
			name: "can find photos with invalid params",
			params: map[string][]string{
				"invalid param 1": {"1"},
				"invalid param 2": {"2"},
			},
			expectedCount: 5000,
		},
		{
			name: "can find photos by IDs",
			params: map[string][]string{
				"id": {"1", "2", "3", "4", "5", "not an ID"},
			},
			expectedCount: 5,
		},
		{
			name: "can find photos by IDs and invalid params",
			params: map[string][]string{
				"id":              {"1", "2", "3", "4", "5"},
				"invalid param 1": {"1"},
			},
			expectedCount: 5,
		},
		{
			name: "can find photos by album IDs",
			params: map[string][]string{
				"albumId": {"1", "5", "not an album ID"},
			},
			expectedCount: 100,
		},
		{
			name: "can find photos by titles",
			params: map[string][]string{
				"title": {
					"accusamus beatae ad facilis cum similique qui sunt",
					"reprehenderit est deserunt velit ipsam",
					"consequuntur consequatur nesciunt vitae deleniti",
					"not a title",
				},
			},
			expectedCount: 3,
		},
		{
			name: "can find photos by URLs",
			params: map[string][]string{
				"url": {
					"https://via.placeholder.com/600/92c952",
					"https://via.placeholder.com/600/771796",
					"https://via.placeholder.com/600/cb0f89",
					"not a URL",
				},
			},
			expectedCount: 3,
		},
		{
			name: "can find photos by thumbnail URLs",
			params: map[string][]string{
				"thumbnailUrl": {
					"https://via.placeholder.com/150/92c952",
					"https://via.placeholder.com/150/771796",
					"https://via.placeholder.com/150/cb0f89",
					"not a thumbnail URL",
				},
			},
			expectedCount: 3,
		},
		{
			name: "can find photos by IDs and album IDs",
			params: map[string][]string{
				"id":      {"1", "2", "3", "4", "5"},
				"albumId": {"1", "2"},
			},
			expectedCount: 5,
		},
		{
			name: "can find photos by IDs and URLs",
			params: map[string][]string{
				"id": {"1", "2", "3", "4", "5"},
				"url": {
					"https://via.placeholder.com/600/92c952",
					"https://via.placeholder.com/600/771796",
				},
			},
			expectedCount: 2,
		},
		{
			name: "can find photos by IDs and thumbnail URLs",
			params: map[string][]string{
				"id": {"1", "2", "3", "4", "5"},
				"thumbnailUrl": {
					"https://via.placeholder.com/150/92c952",
					"https://via.placeholder.com/150/771796",
				},
			},
			expectedCount: 2,
		},
		{
			name: "can find photos by album IDs and URLs",
			params: map[string][]string{
				"albumId": {"1", "2"},
				"url": {
					"https://via.placeholder.com/600/92c952",
					"https://via.placeholder.com/600/771796",
				},
			},
			expectedCount: 2,
		},
		{
			name: "can find photos by album IDs and thumbnail URLs",
			params: map[string][]string{
				"albumId": {"1", "2"},
				"thumbnailUrl": {
					"https://via.placeholder.com/150/92c952",
					"https://via.placeholder.com/150/771796",
				},
			},
			expectedCount: 2,
		},
		{
			name: "can find photos by URLs and thumbnail URLs",
			params: map[string][]string{
				"url": {
					"https://via.placeholder.com/600/92c952",
					"https://via.placeholder.com/600/771796",
				},
				"thumbnailUrl": {
					"https://via.placeholder.com/150/92c952",
					"https://via.placeholder.com/150/771796",
				},
			},
			expectedCount: 2,
		},
		{
			name: "returns no results when no photos meet filter criteria",
			params: map[string][]string{
				"id":      {"1", "2", "3", "4", "5"},
				"albumId": {"99", "100"},
			},
			expectedCount: 0,
		},
	}

	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			client := photos.NewClient(http.DefaultClient)
			ps, err := client.Photos(context.Background(), testCase.params)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if len(ps) != testCase.expectedCount {
				t.Errorf("Expected %d photos, got %d", testCase.expectedCount, len(ps))
			}
		})
	}
}
