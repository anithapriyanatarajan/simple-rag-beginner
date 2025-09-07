# Multi-stage build for rag-api
FROM golang:1.23-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN cd cmd/rag-api && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app ./
COPY web ./web
RUN adduser -D -s /bin/sh appuser
RUN chown -R appuser:appuser ./web
USER appuser
EXPOSE 8080
CMD ["./app"]
