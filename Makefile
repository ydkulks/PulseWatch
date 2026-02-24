build-proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/v1/pulsewatch/pulsewatch.proto

lint:
	go fmt ./...

build-linux:
	go build -o bin/pulsewatchserver cmd/server/main.go

build-macos:
	GOOS=darwin GOARCH=amd64 go build -o bin/pulsewatchserver_x64 cmd/server/main.go

build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/pulsewatchserver_x64.exe cmd/server/main.go

build: build-linux build-macos build-windows

run:
	./bin/pulsewatch
