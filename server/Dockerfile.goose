FROM golang:1.24-alpine AS build
RUN apk add --no-cache alpine-sdk

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o main cmd/goose/main.go

FROM alpine:latest AS migrations
WORKDIR /app
COPY --from=build /app/main /app/main
CMD ["./main"]
