# Chat Service

**Chat Service** is a real-time messaging microservice built with Go and gRPC. It supports chat creation, real-time communication via bi-directional streaming, and secure message delivery.

---

## 🚀 Quick Start

### Requirements

- Go 1.19+
- Docker & Docker Compose

### 1. Clone the project

```bash
git clone https://github.com/your-org/chat-service.git
cd chat-service
````

### 2. Setup and run app(recommended)

```bash
make setup
```

This will:

- Install all dev dependencies
- Generate gRPC + gateway code
- Start services via Docker Compose
- Run app

---

## 📦 Features

- Create / Get / Delete chats
- Real-time bi-directional gRPC streaming
- Secure message delivery with persistence
- Auth integration with access control
- Protobuf + Swagger + Gateway generation

---

## 🧰 Makefile Commands

| Command                  | Description                                                 |
| ------------------------ | ----------------------------------------------------------- |
| `make setup`             | Install deps, generate code, and start services             |
| `make install-deps`      | Install CLI tools (protoc, gateway, validate, linter)       |
| `make generate`          | Generate gRPC, gateway, validators, Swagger + static assets |
| `make generate-chat-api` | Generate code from `chat.proto`                             |
| `make vendor-proto`      | Download external `.proto` definitions                      |
| `make up`                | Start PostgreSQL, Redis, and service using Docker Compose   |
| `make lint`              | Run static analysis via golangci-lint                       |
| `make test`              | Run unit tests with retries and coverage                    |
| `make test-coverage`     | Generate coverage report and open it in browser             |
| `make gen-cert`          | Create local TLS certificates for gRPC testing              |
| `make grpc-load-test`    | Run performance/load test with TLS using `ghz`              |

---

## 🛠 Tech Stack

- Go 1.19+
- gRPC + Protobuf
- PostgreSQL
- Redis
- Docker Compose
- OpenTracing
- Zap (structured logging)

---

## 📊 Observability

### Logging

- [Zap](https://github.com/uber-go/zap) for structured logging
- Includes tracing, user IDs, request context

### Tracing

- [OpenTracing](https://opentracing.io/) spans are injected into gRPC and business logic

---

## 🔐 Security & Auth

- All RPC calls require a valid token
- Integration with Auth Service
- Token validation + role-based access control
- TLS certificate support (`make gen-cert`)

---

## 📁 Project Layout

```
├── api/                      # Protobuf definitions
├── cmd/                      # Entrypoints (main.go)
├── internal/                 # Private application code
│   ├── app/                  # Application wiring (DI, lifecycle)
│   ├── client/               # External service clients (e.g., Auth)
│   ├── config/               # Config loading
│   ├── converter/            # Data transformers between layers
│   ├── interceptor/          # gRPC interceptors (auth, logging, etc.)
│   ├── model/                # Domain models and constants
│   ├── ratelimiter/          # Rate limiting logic
│   ├── repository/           # Storage access (Postgres, etc.)
│   ├── service/              # Business logic
│   └── transport/
│       └── handlers/         # gRPC handlers
├── migrations/               # Database schema (Goose)
├── pkg/                      # Generated code and shared helpers
├── vendor.protogen/          # External proto dependencies
├── Makefile                  # Dev utility commands
├── local.env.example         # Example environment variables


```

---

## 📎 Example: Message Flow

1. **Connect:** Client starts stream via `ConnectChat` (bi-directional)
2. **Send Message:** Client sends via `SendMessage`
3. **Server:** Validates, saves, broadcasts to all connected participants
4. **Receive:** Clients receive messages in real-time from stream

---

## 🔗 Dependencies

- **Auth Service** – JWT validation, user identity
- **PostgreSQL** – Persistent storage
- **Redis** – (Optional) for future scalability

---

## License

MIT
