gen:
	protoc ./proto/*.proto  -I=./proto --go_out=.  --go-grpc_out=.

clean:
	rm pb/*

server:
	go run cmd/server/main.go -port 8080

client:
	go run cmd/client/main.go -address 127.0.0.1:8080


.PHONY: clean gen server client 