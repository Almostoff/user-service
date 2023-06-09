FROM golang:1.18.3-alpine3.15 as gobuild

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .

RUN go build -o api ./cmd/api

FROM alpine:latest
WORKDIR /app

COPY --from=gobuild /app .
CMD ["./api"]