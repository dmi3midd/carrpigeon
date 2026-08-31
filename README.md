<p align="center">
  <img src="assets/logo.png" alt="carrpigeo logo" width="300">
</p>

# Carrpigeo

Carrpigeo is a lightweight email sending API built with Go. It sends emails via SMTP and stores every sent message in a PostgreSQL database.

## Prerequisites

- [Go](https://go.dev/dl/) 1.26+
- [PostgreSQL](https://www.postgresql.org/download/) 17+ (local install **or** Docker)
- [Docker & Docker Compose](https://docs.docker.com/get-docker/) (optional — for containerized setup)
- SMTP credentials (e.g. Gmail [App Password](https://support.google.com/accounts/answer/185833))

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/dmi3midd/carrpigeo.git
cd carrpigeo
```

### 2. Run the setup script

This creates the `storage/` directory, copies `config.example.yaml` → `config.yaml`, and creates the log file:

```bash
make setup
```

### 3. Configure

Edit `config.yaml` with your credentials:

```yaml
# Configure it for you.
database:
  name: carrpigeo
  host: postgres
  port: 5432
  user: carrpigeo_user
  password: carrpigeo_9876
  sslmode: disable
  maxOpenConns: 10
  maxIdleConns: 10
  maxIdleTime: 15m

httpServer:
  address: 0.0.0.0:2500

smtp:
  host: smtp.gmail.com
  port: 587
  user: your-email@gmail.com
  password: your-app-password

log:
  logPath: ./storage/carrpigeo.log
```

> **Note:** For Gmail, you need to generate an [App Password](https://support.google.com/accounts/answer/185833) — regular account passwords won't work with SMTP.

## Running the Project

### Option A: Local (without Docker)

> **Note** Make sure PostgreSQL is running locally and the database exist.

Run:

```bash
make run
```

Or use `air` for auto-reloading:

```bash
make watch
```

The server starts on `http://localhost:2500`.

### Option B: Docker (App + PostgreSQL)

Make sure `host: postgres` is configured in `config.yaml`, then:

```bash
make docker-run
```

This builds the Docker image and starts both the application and PostgreSQL containers in the background. Database migrations run automatically on startup.

## API Endpoints

### `GET /health`

Returns database health status.

**Response:** `200 OK`

```json
{
    "idle": "1",
    "in_use": "0",
    "max_idle_closed": "0",
    "max_lifetime_closed": "0",
    "message": "It's healthy",
    "open_connections": "1",
    "status": "up",
    "wait_count": "0",
    "wait_duration": "0s"
}
```

---

### `POST /send/email`

Sends an email and saves it to the database.

**Request body:**

```json
{
  "to": "recipient@example.com",
  "subject": "Hello",
  "body": "This is the email body"
}
```

**Response:** `202 Accepted`

**Error responses:**

| Code | Description |
|------|-------------|
| 400  | Invalid request body |
| 500  | Failed to send email or save to database |

---

### `POST /templates`

Uploads and registers a new HTML template.

**Request:** `multipart/form-data`

| Parameter | Type | Description |
|-----------|------|-------------|
| `name`    | text | The name of the HTML template |
| `file`    | file | HTML file containing template placeholders (e.g. `{{.Name}}`) |

**Response:** `202 Accepted`

```json
{
  "id": "cno3c9u8u1401p0q7sdg"
}
```

---

### `DELETE /templates`

Deletes a registered HTML template by ID.

**Query Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id`  | string | Yes | The ID of the HTML template to delete |

**Response:** `200 OK`

---

### `POST /send/email/template`

Sends an email using a registered HTML template and performs dynamic data insertion.

**Request body:**

```json
{
  "to": "recipient@example.com",
  "subject": "Welcome!",
  "template_id": "cno3c9u8u1401p0q7sdg",
  "data": {
    "Name": "John Doe",
    "PromoCode": "WELCOME2026"
  }
}
```

**Response:** `202 Accepted`

#### Template Syntax & Dynamic Variables

Inside your HTML template files, you can define placeholders using Go's `html/template` syntax. Placeholders must be prefixed with a dot `.` followed by the field/key name, wrapped in double curly braces: `{{.FieldName}}`.

For example, if you upload a template file containing:
```html
<h1>Hello, {{.Name}}!</h1>
<p>Your promotional code is: <strong>{{.PromoCode}}</strong></p>
```

You can populate these values by sending them in the `data` field of the `POST /send/email/template` request:
```json
{
  "to": "recipient@example.com",
  "subject": "Your Promo Code",
  "template_id": "cno3c9u8u1401p0q7sdg",
  "data": {
    "Name": "John Doe",
    "PromoCode": "WELCOME2026"
  }
}
```
*Note: The keys in the JSON `data` object must match the placeholder names exactly (case-sensitive).*

**Error responses:**

| Code | Description |
|------|-------------|
| 400  | Invalid request body or missing required fields |
| 500  | Failed to retrieve template, parse/execute template, send email, or save to database |

## Make Targets

| Command | Description |
|---------|-------------|
| `make setup` | Initialize project (directories and config file) |
| `make run` | Run the application locally |
| `make build` | Compile the binary locally |
| `make watch` | Live reload with [air](https://github.com/air-verse/air) |
| `make test` | Run tests |
| `make tidy` | Run go mod tidy |
| `make clean` | Remove the compiled binary |
| `make docker-build` | Build Docker images |
| `make docker-run` | Build and start all containers in background |
| `make docker-down` | Stop all Docker containers |
| `make docker-down-v` | Stop containers and remove volumes (clean database reset) |
| `make docker-logs` | Follow app container logs |
| `make docker-restart` | Restart application container |

## Project Structure

## License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.
