package main

// Testing the static file server
import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStaticFileServer(t *testing.T) {
	r := setupRouter()
	mockServer := httptest.NewServer(r)

	// We want to hit the `GET /public/` route to get the index.html file response
	resp, err := http.Get(mockServer.URL + "/public/")
	if err != nil {
		t.Fatal(err)
	}

	// We want our status to be 200 (ok)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status should be 200, got %d", resp.StatusCode)
	}
	// Test the content-type header is "text/html; charset=utf-8"
	// so that we know that an html file has been served
	contentType := resp.Header.Get("Content-Type")
	expectedContentType := "text/html; charset=utf-8"

	if expectedContentType != contentType {
		t.Errorf("Wrong content type, expected %s, got %s", expectedContentType, contentType)
	}

}
