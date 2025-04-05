package wrapper

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var (
	SearchApiBaseUrl = "https://api.github.com/search"
)

type GithubSearchCodeResult struct {
	FileURL string
	Repo    string
}

type GithubSearchApiClient struct {
	Token   string
	Client  *http.Client
	BaseURL string
}

func NewGithubSearchApiClient(token string) *GithubSearchApiClient {
	return &GithubSearchApiClient{
		Token:   token,
		Client:  &http.Client{},
		BaseURL: SearchApiBaseUrl,
	}
}

func (g *GithubSearchApiClient) SearchCode(ctx context.Context, term, user string) ([]GithubSearchCodeResult, error) {
	query := term
	if user != "" {
		query += " user:" + user
	}
	apiURL := fmt.Sprintf(g.BaseURL+"/code?q=%s", url.QueryEscape(query))

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+g.Token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := g.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	var jsonResp struct {
		Items []struct {
			HTMLURL    string `json:"html_url"`
			Repository struct {
				FullName string `json:"full_name"`
			} `json:"repository"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&jsonResp); err != nil {
		return nil, err
	}

	results := []GithubSearchCodeResult{}
	for _, item := range jsonResp.Items {
		results = append(results, GithubSearchCodeResult{
			FileURL: item.HTMLURL,
			Repo:    item.Repository.FullName,
		})
	}

	return results, nil
}
