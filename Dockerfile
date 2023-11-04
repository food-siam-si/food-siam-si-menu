FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o server main.go

FROM alpine:latest as runner

WORKDIR /app

COPY --from=builder /app/server .

ENV GO_ENV production

EXPOSE 8080

CMD ["./server"]
