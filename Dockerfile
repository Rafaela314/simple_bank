FROM golang:1.26.1-alpine3.23 AS builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY . .
RUN go build -o main main.go

FROM alpine:3.23
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/main /app/main

COPY db/migration ./migration

EXPOSE 8080
CMD ["/app/main"]