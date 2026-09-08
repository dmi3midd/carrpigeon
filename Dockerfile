# Stage 1: Build Web UI
FROM node:22-alpine AS ui-builder

WORKDIR /app/webui

COPY webui/package*.json ./
RUN npm install

COPY webui/ ./
RUN npm run build

# Stage 2: Build Go binary
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/carrpigeon ./cmd/api/main.go

# Stage 3: Final minimal runtime image
FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/carrpigeon .
COPY --from=ui-builder /app/webui/dist ./webui/dist

RUN mkdir -p /app/storage

EXPOSE 2500

ENTRYPOINT ["./carrpigeon"]
