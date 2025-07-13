FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go get -d ./...

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/server-binary ./server/*.go

# 2nd stage runtime image
FROM alpine:3.18

WORKDIR /app

ENV ENV=production

COPY --from=builder /app/server-binary .
COPY --from=builder /app/database/migrations /app/database/migrations

EXPOSE 8080

CMD ["./server-binary"]