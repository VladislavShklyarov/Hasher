FROM golang:1.23 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

CMD ["./main"]
