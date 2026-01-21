FROM golang:1.24-alpine AS modules
WORKDIR /modules
COPY go.mod go.sum ./
RUN go mod download

FROM golang:1.24-alpine AS builder
COPY --from=modules /go/pkg /go/pkg
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o ./bin/main ./cmd/gateway/main.go

FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /build/bin/main .
EXPOSE 8080
CMD ["./main"]