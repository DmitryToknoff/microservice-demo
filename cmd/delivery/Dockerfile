FROM golang:1.24.1-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/delivery-exe ./cmd/delivery

FROM alpine:3.23 AS runner
WORKDIR /app
COPY --from=builder /app/delivery-exe /app/delivery-exe
COPY --from=builder /app/migrations /app/migrations
ENV PROJECT_ROOT=/app
CMD ["/app/delivery-exe"]