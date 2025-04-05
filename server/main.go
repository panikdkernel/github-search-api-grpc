package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"

	"github.com/panikdkernel/github-search-api-grpc/config"
	pb "github.com/panikdkernel/github-search-api-grpc/githubsearch_proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGithubSearchServiceServer
}

func (s *server) Search(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	log.Printf("Search called with term: %s, user: %s", req.SearchTerm, req.User)

	// Prepare GitHub Search API URL
	baseURL := config.GithubSearchApiCodeBaseUrl
	query := req.SearchTerm
	if req.User != "" {
		query += " user:" + req.User
	}
	apiURL := fmt.Sprintf("%s?q=%s", baseURL, url.QueryEscape(query))

	// Create HTTP client and request
	client := &http.Client{}
	httpReq, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	// GitHub token is required
	token := config.GithubToken

	httpReq.Header.Set("Authorization", "Bearer "+token)

	// Make the request
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse response
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

	// Map to gRPC response format
	results := []*pb.Result{}
	for _, item := range jsonResp.Items {
		results = append(results, &pb.Result{
			FileUrl: item.HTMLURL,
			Repo:    item.Repository.FullName,
		})
	}

	return &pb.SearchResponse{Results: results}, nil
}

func main() {
	grpc_port := config.Port
	lis, err := net.Listen("tcp", ":"+grpc_port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	// GitHub token is required
	token := config.GithubToken
	if token == "" {
		log.Fatal("missing GitHub token: set GITHUB_TOKEN environment variable")
	}
	pb.RegisterGithubSearchServiceServer(s, &server{})

	log.Println("gRPC server running on :" + grpc_port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
