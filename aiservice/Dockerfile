FROM golang:1.25-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o aiservice ./aiservice

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /build/aiservice .
COPY aiservice/etc/ai.yaml etc/ai.yaml
EXPOSE 8081
ENTRYPOINT ["./aiservice", "-f", "etc/ai.yaml"]
