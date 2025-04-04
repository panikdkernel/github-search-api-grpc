build_proto:
	rm -rf githubsearch_proto && mkdir -p githubsearch_proto
	protoc --proto_path=proto/ --go_out=./githubsearch_proto \
	--go-grpc_out=./githubsearch_proto \
	--go_opt=paths=source_relative \
	--go-grpc_opt=paths=source_relative \
	proto/github_search.proto