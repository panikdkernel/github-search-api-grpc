# Use official Go image as base
FROM golang:1.23.3 as builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Install protoc and plugins
RUN apt-get update && \
    apt-get install -y unzip curl && \
    curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v21.12/protoc-21.12-linux-x86_64.zip && \
    unzip protoc-21.12-linux-x86_64.zip -d /usr/local && \
    rm -rf protoc-21.12-linux-x86_64.zip

# install required plugins
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

ENV PATH="/go/bin:$PATH"

# Generate Proto code
RUN make build_proto

# Build the Go binary
RUN go build -o server ./server

# Final lightweight image
FROM gcr.io/distroless/base-debian11

WORKDIR /app

# Copy built binary from builder
COPY --from=builder /app/server .


RUN apt update && apt install libc6
# Set required environment variables here (or set at runtime)
# ENV GITHUB_TOKEN=your_token
# ENV GRPC_PORT=9001

# Run the binary
ENTRYPOINT ["/app/server"]
