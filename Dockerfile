FROM golang:1.20-alpine AS builder

WORKDIR /app
COPY . ./

RUN go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -o /hydroponics-metrics-collector ./cmd/api/main.go

FROM alpine:3.18

COPY --from=builder /hydroponics-metrics-collector /hydroponics-metrics-collector

EXPOSE 8080
CMD ["/hydroponics-metrics-collector"]
