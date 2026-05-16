FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . . 
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/gateway

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/configs/config.yaml ./configs/
COPY --from=builder /app/lua/token_bucket.lua ./lua/
COPY --from=builder /app/.env .
EXPOSE 8080
CMD ["./main"]
