build-proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/v1/pulsewatch/services.proto

lint:
	go fmt ./...
