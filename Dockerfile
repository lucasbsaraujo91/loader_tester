# Etapa de build
FROM golang:1.23 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go build -o load_tester .

# Etapa final
FROM debian:bookworm-slim

WORKDIR /app
COPY --from=builder /app/load_tester /usr/local/bin/

ENTRYPOINT ["load_tester"]
