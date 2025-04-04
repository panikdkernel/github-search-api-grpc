package main

import (
	"context"
	"log"
	"net"

	pb "github.com/panikdkernel/github-search-api-grpc/githubsearch_proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGithubSearchServiceServer
}

func (s *server) Search(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	// Dummy implementation
	log.Printf("Search called with term: %s, user: %s", req.SearchTerm, req.User)

	// Sample response
	results := []*pb.Result{
		{FileUrl: "https://github.com/example/file1", Repo: "example/repo1"},
		{FileUrl: "https://github.com/example/file2", Repo: "example/repo2"},
	}
	return &pb.SearchResponse{Results: results}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGithubSearchServiceServer(s, &server{})

	log.Println("gRPC server running on :9001")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
