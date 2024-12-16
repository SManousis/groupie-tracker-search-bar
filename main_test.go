package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func mockFetchArtists(artists *[]Artist) error {
	*artists = []Artist{
		{
			ID:           1,
			Name:         "Test Artist",
			Image:        "test_image.jpg",
			Members:      []string{"Member 1", "Member 2"},
			CreationDate: 2000,
			FirstAlbum:   "2001",
			Locations:    []string{"Location 1", "Location 2"},
			Dates:        []string{"2023-01-01", "2023-01-02"},
			Relation: map[string]string{
				"2023-01-01": "Location 1",
				"2023-01-02": "Location 2",
			},
		},
	}
	return nil
}

func TestServeHome(t *testing.T) {
	// Mock FetchArtists
	originalFetchArtists := FetchArtists
	defer func() { FetchArtists = originalFetchArtists }()
	FetchArtists = mockFetchArtists

	// Create the test request
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	// Call the handler
	serveHome(w, req)

	// Get the response
	res := w.Result()
	defer res.Body.Close()

	// Check the HTTP status code
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", res.StatusCode)
	}

	// Get the response body as a string
	body := w.Body.String()

	// Look for the JavaScript content
	jsSnippet := `window.artistData = [`
	if !strings.Contains(body, jsSnippet) {
		t.Errorf("Expected response body to contain JavaScript 'window.artistData', but it didn't")
	}

	// Look for the specific artist name in the JavaScript content
	expectedArtist := `name: "Test Artist"`
	if !strings.Contains(body, expectedArtist) {
		t.Errorf("Expected response body to contain %q, but it didn't. Full response:\n%v", expectedArtist, body)
	}
}

func TestServeHomeNotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/non-existent-path", nil)
	w := httptest.NewRecorder()

	serveHome(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status Not Found, got %v", res.StatusCode)
	}
}
