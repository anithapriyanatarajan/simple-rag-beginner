# Multi-stage build for embedder
FROM golang:1.23-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN cd cmd/embedder && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app ./
RUN adduser -D -s /bin/sh appuser
USER appuser
EXPOSE 8080
CMD ["./app"]
