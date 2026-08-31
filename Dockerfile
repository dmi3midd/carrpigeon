FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/carrpigeon ./cmd/api/main.go

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/carrpigeon .

RUN mkdir -p /app/storage

EXPOSE 2500

ENTRYPOINT ["./carrpigeon"]
