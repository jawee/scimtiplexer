FROM golang:1.24-alpine AS build
RUN apk add --no-cache alpine-sdk

FROM build AS app

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o main cmd/api/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=app /app/main .
COPY --from=app /app/.env .

EXPOSE 8080

CMD ["./main"]
