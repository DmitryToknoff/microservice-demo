FROM golang:1.25-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/gateway-exe ./cmd/gateway

FROM alpine:3.23 AS runner
WORKDIR /app
COPY --from=builder /app/gateway-exe /app/gateway-exe
COPY --from=builder /app/migrations /app/migrations
ENV PROJECT_ROOT=/app
CMD ["/app/gateway-exe"]