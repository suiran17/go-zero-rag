FROM golang:1.25-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o docservice ./docservice

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /build/docservice .
COPY docservice/etc/doc.yaml etc/doc.yaml
EXPOSE 8080
ENTRYPOINT ["./docservice", "-f", "etc/doc.yaml"]
