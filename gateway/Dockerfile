FROM golang:1.25-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o gateway ./gateway

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /build/gateway .
COPY gateway/etc/gateway.yaml etc/gateway.yaml
EXPOSE 8888
ENTRYPOINT ["./gateway", "-f", "etc/gateway.yaml"]
