FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o likho cmd/likho/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/likho /app/likho
COPY templates /app/templates
COPY config.yaml /app/config.yaml
COPY assets /app/assets

ENTRYPOINT ["/app/likho"]
