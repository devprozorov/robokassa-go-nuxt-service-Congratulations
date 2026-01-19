# syntax=docker/dockerfile:1
FROM golang:1.22-alpine AS build
WORKDIR /src
RUN apk add --no-cache ca-certificates tzdata && update-ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/happy-api ./cmd/api

FROM alpine:3.20
WORKDIR /app
RUN adduser -D -H -u 10001 appuser && apk add --no-cache ca-certificates tzdata && update-ca-certificates
COPY --from=build /out/happy-api /app/happy-api
USER appuser
EXPOSE 8080
ENTRYPOINT ["./happy-api"]
