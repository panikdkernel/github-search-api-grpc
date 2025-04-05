# GitHub Search API - gRPC Service

A Go-based gRPC microservice that uses the GitHub Search API to search for code files, optionally filtered by GitHub username. It returns file URLs and their repositories.

---

## ✨ Features

- gRPC service using Protocol Buffers
- GitHub Search API integration
- Optional username filter
- Supports local and Docker-based deployment

---

## ⚙️ Environment Variables

| Variable       | Description                         | Required |
|----------------|-------------------------------------|----------|
| `GITHUB_TOKEN` | GitHub Personal Access Token        | ✅       |
| `GRPC_PORT`    | gRPC Server Port (default: `9001`)  | ❌       |

> **Note**: The GitHub token is required to avoid API rate limits and is now mandatory for GitHub Search API requests.

---

## 🧪 Run Locally

### 1. Clone the repo

```bash
git clone https://github.com/YOUR_USERNAME/github-search-api-grpc.git
cd github-search-api-grpc
```

### 2. Set environment variables

```bash
export GITHUB_TOKEN=your_github_token
export GRPC_PORT=9001
```

### 3. Generate Go code from .proto

```bash
make build_proto
```

### 4. Run the server

```bash
go run server/main.go
```

---

## 🐳 Run with Docker

### 1. Build the Docker image

```bash
docker build -t github-search-api-grpc .
```

### 2. Run the container

```bash
docker run -e GITHUB_TOKEN=your_github_token -e GRPC_PORT=9001 -p 9001:9001 github-search-api-grpc
```

---

## 🔌 Test with grpcurl

```bash
grpcurl -plaintext \
  -proto proto/github_search.proto \
  -d '{"search_term": "docker", "user": "torvalds"}' \
  localhost:9001 \
  githubsearch.GithubSearchService/Search
```
---

## 📜 License

MIT

---

## 👨‍💻 Author

Onkar Singh – [@panikdkernel](https://github.com/panikdkernel)