# Use buildx syntax for multi-platform builds
FROM --platform=$BUILDPLATFORM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .
ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o likho cmd/likho/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/likho /app/likho
COPY config.yaml /app/config.yaml
COPY themes /app/themes

ENTRYPOINT ["/app/likho"]
