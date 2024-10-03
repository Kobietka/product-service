FROM golang:1.23.1-alpine3.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /app/main /app/cmd/service/main.go

FROM scratch

COPY --from=builder /app/main /main

ENTRYPOINT ["/main"]