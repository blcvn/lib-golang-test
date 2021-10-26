
all: client server

protoc:
	@echo "Generating Go files"
	cd proto && protoc --go_out=. --go-grpc_out=. \
		--go-grpc_opt=paths=source_relative --go_opt=paths=source_relative *.proto

server: 
	@echo "Running server"
	go run server.go 

client: protoc
	@echo "Building client"
	go build -o client \
		github.com/pahanini/go-grpc-bidirectional-streaming-example/src/client

clean:
	go clean github.com/pahanini/go-grpc-bidirectional-streaming-example/...
	rm -f server client

.PHONY: client server protoc
