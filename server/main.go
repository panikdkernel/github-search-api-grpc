package main

import (
	"context"
	"log"
	"net"

	"github.com/panikdkernel/github-search-api-grpc/config"
	pb "github.com/panikdkernel/github-search-api-grpc/githubsearch_proto"
	"github.com/panikdkernel/github-search-api-grpc/internal/wrapper"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGithubSearchServiceServer
}

func (s *server) Search(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	log.Printf("Search called with term: %s, user: %s", req.SearchTerm, req.User)

	githubSearchClient := wrapper.NewGithubSearchApiClient(config.GithubToken)

	// Make the request
	results, err := githubSearchClient.SearchCode(ctx, req.SearchTerm, req.User)
	if err != nil {
		return nil, err
	}

	// Map to gRPC response format
	var grpcResults []*pb.Result
	for _, item := range results {
		grpcResults = append(grpcResults, &pb.Result{
			FileUrl: item.FileURL,
			Repo:    item.Repo,
		})
	}

	return &pb.SearchResponse{Results: grpcResults}, nil
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
