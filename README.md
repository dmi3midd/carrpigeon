<p align="center">
  <img src="assets/logo.png" alt="carrpigeon logo" width="300">
</p>

# Carrpigeon

**Carrpigeon** is a lightweight, email notification service built in Go. It features an asynchronous background email worker with retry policies, sharded in-memory caching, full PostgreSQL persistence, a built-in terminal-styled Web UI, and interactive Swagger API documentation.

---

## Key Features

- **Asynchronous Worker Pool**: Non-blocking email dispatching with configurable worker pools, queue polling intervals, and retry attempts.
- **Embedded Web UI**: Built-in CLI/terminal-styled web interface (React 19 + TypeScript + Zustand) served directly by the Go binary at `http://localhost:2500/`.
- **Dynamic Template Engine**: Support for HTML and plaintext templates using Go's `html/template` and `text/template` syntax with automatic placeholder detection like `{{.Name}}`.
- **Recipients & Group Management**: Directory for individual email recipients and distribution groups with dynamic membership.
- **Interactive Swagger Documentation**: Full OpenAPI/Swagger UI available out-of-the-box at `/swagger/index.html`.

---

## Prerequisites

Before running the project, ensure you have:

- **[Docker & Docker Compose](https://docs.docker.com/get-docker/)** (Recommended — covers Go, PostgreSQL, and Web UI build in containers).
- **[Go](https://go.dev/dl/)** 1.24+ (if running locally without Docker).
- **[PostgreSQL](https://www.postgresql.org/download/)** 17+ (if running database locally outside Docker).
- **[Node.js](https://nodejs.org/)** 20+ & **npm** (only needed if developing or building the Web UI locally).
- **SMTP Credentials**

---

## Quick Start (Docker Compose — Recommended)

The simplest way to start the complete stack (Carrpigeon API + Web UI + PostgreSQL) is using Docker Compose:

### 1. Clone the repository

```bash
git clone https://github.com/dmi3midd/carrpigeon.git
cd carrpigeon
```

### 2. Initialize project configuration

Run the setup script, which creates the `./storage` directory and copies `config.example.yaml` to `config.yaml`:

```bash
make setup
```

### 3. Configure credentials

Open `config.yaml` and configure your SMTP credentials:

```yaml
email:
  smtp:
    host: smtp.gmail.com
    port: 587
    user: your-email@gmail.com
    password: your-app-password
```

> **Note for Gmail users:** Standard Google account passwords will not work. You must generate a 16-character [App Password](https://support.google.com/accounts/answer/185833) with 2-Step Verification enabled.

Ensure that `postgres.host` is set to `postgres` (the Docker service name).

### 4. Start services

```bash
make docker-run
```

*Or via Docker Compose directly:*

```bash
docker compose up --build -d
```

This will:

1. Build the Web UI production bundle with Node.js.
2. Compile the Go backend binary.
3. Start the PostgreSQL container with health checks.
4. Run all database migrations automatically.
5. Launch Carrpigeon on `http://localhost:2500`.

### 5. Access the Service

- **Web UI Console:** [http://localhost:2500](http://localhost:2500)
- **Swagger Documentation:** [http://localhost:2500/swagger/index.html](http://localhost:2500/swagger/index.html)
- **Health Check Endpoint:** [http://localhost:2500/health](http://localhost:2500/health)

---

## API Reference

Interactive Swagger documentation is accessible at **`http://localhost:2500/swagger/index.html`**.

---

## Template Syntax

Templates use Go's standard `html/template` or `text/template` syntax. Placeholders must be wrapped in `{{.FieldName}}`:

```html
<h2>Welcome aboard, {{.Name}}!</h2>
<p>Your activation link is: <a href="{{.Link}}">{{.Link}}</a></p>
```

---

## Makefile Command Reference

| Command | Description |
| --- | --- |
| `make setup` | Initialize storage directory and copy `config.example.yaml` → `config.yaml` |
| `make docker-run` | Build and launch all containers (API + Web UI + PostgreSQL) in the background |
| `make docker-down` | Stop and remove running Docker containers |
| `make docker-down-v` | Stop containers and remove volumes (**resets PostgreSQL database**) |
| `make docker-logs` | Stream logs from the application container |
| `make docker-restart` | Restart the application container |
| `make run` | Run Go backend locally on host machine |
| `make watch` | Run Go backend locally with live-reload ([Air](https://github.com/air-verse/air)) |
| `make build` | Compile Go binary locally (`./main`) |
| `make test` | Run all backend Go unit and integration tests |
| `make ui-dev` | Start Web UI development server with Vite hot reload (`http://localhost:5173`) |
| `make ui-build` | Build Web UI production bundle to `webui/dist` |
| `make ui-test` | Run Web UI unit tests (Vitest) |
| `make swagger` | Regenerate Swagger/OpenAPI documentation (`swag init`) |
| `make tidy` | Tidy and verify Go module dependencies (`go mod tidy`) |
| `make clean` | Remove compiled local binary |
