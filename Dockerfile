FROM golang:1.22.2 AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

COPY .env .

RUN CGO_ENABLED=0 GOOS=linux go build -C ./cmd -a -installsuffix cgo -o /app/myapp

FROM alpine:latest

WORKDIR /app

RUN mkdir -p ./logs

COPY --from=builder /app/myapp .
COPY --from=builder /app/logs/app.log ./logs/

COPY --from=builder /app/.env .

EXPOSE 8081

CMD ["./myapp"]
