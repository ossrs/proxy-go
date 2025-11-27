# Build stage
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache make

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make build

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/srs-proxy .

CMD ["./srs-proxy"]