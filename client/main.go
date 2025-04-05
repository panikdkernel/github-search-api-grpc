package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/panikdkernel/github-search-api-grpc/githubsearch_proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:9001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to connect to the grpc server")
	}

	defer conn.Close()

	c := pb.NewGithubSearchServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()
	search_result, err := c.Search(ctx, &pb.SearchRequest{})
	if err != nil {
		log.Fatal("failed to search: ", err)
	}
	fmt.Print(search_result)
}
