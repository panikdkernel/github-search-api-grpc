package wrapper

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSearchCode_Success(t *testing.T) {
	mockJSON := `{
		"items": [
			{
				"html_url": "https://github.com/example/repo1/file.go",
				"repository": {
					"full_name": "example/repo1"
				}
			}
		]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockJSON))
	}))
	defer server.Close()

	client := &GithubSearchApiClient{
		Token:   "test-token",
		Client:  server.Client(),
		BaseURL: server.URL,
	}

	results, err := client.SearchCode(context.Background(), "test", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	if results[0].FileURL != "https://github.com/example/repo1/file.go" {
		t.Errorf("unexpected file URL: %s", results[0].FileURL)
	}
	if results[0].Repo != "example/repo1" {
		t.Errorf("unexpected repo: %s", results[0].Repo)
	}
}

func TestSearchCode_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	client := &GithubSearchApiClient{
		Token:   "invalid-token",
		Client:  server.Client(),
		BaseURL: server.URL,
	}

	_, err := client.SearchCode(context.Background(), "test", "")
	if err == nil || !strings.Contains(err.Error(), "GitHub API error") {
		t.Errorf("expected error containing 'GitHub API error', got %v", err)
	}
}
