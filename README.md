# PulseWatch
PulseWatch is a gRPC-based notification microservice that monitors multiple running programs and sends real-time alerts to various communication channels (Email, SMS, Discord). It helps track the health and status of local or distributed processes through a unified, event-driven system.

## Prerequisits

- Go
- Protobuf compiler
    ```bash
    apt-get install protobuf-compiler
    #or
    pacman -S protobuf
    ```

- Go plugins
    ```bash
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    ```

    - Update the path so that the `protoc` can find the plugins
        ```bash
        export PATH="$PATH:$(go env GOPATH)/bin"
        ```

## Code Architecture

```text
/cmd
   server/main.go
/internal
   /grpc
      notification_service.go   # gRPC handlers
   /domain
      notification.go           # core business logic
      notifier_factory.go
      strategies.go
   /infrastructure
      repository.go
      worker.go
/pkg
   proto/notification.proto

```

## Tech Stack
- Go (backend, concurrency, gRPC server)
- Protocol Buffers (service contracts)
- gRPC Streaming (real-time event updates)
- Docker (containerization)
- Twilio / SendGrid / Slack APIs (notification channels)
- PostgreSQL / Redis (optional) (event persistence)
