gen:
	protoc ./proto/*.proto  -I=./proto --go_out=.  --go-grpc_out=.

clean:
	rm pb/*

build-server:
	go build -o ./bin/server ./cmd/server

build-client:
	go build -o ./bin/client ./cmd/client

server: build-server
	./bin/server -port 8080

client: build-client
	./bin/client -address=127.0.0.1:8080

client-auth: build-client
	./bin/client -address=127.0.0.1:8080 -service=auth

.PHONY: clean gen server client auth-client