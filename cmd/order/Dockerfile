FROM golang:1.24.1-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/order-exe ./cmd/order

FROM alpine:3.23 AS runner
WORKDIR /app
COPY --from=builder /app/order-exe /app/order-exe
COPY --from=builder /app/migrations /app/migrations
ENV PROJECT_ROOT=/app
CMD ["/app/order-exe"]