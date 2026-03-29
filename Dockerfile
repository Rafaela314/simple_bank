FROM golang:1.26.1-alpine3.23 AS builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY . .
RUN go build -o main main.go

FROM alpine:3.23
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./migration

EXPOSE 8080 9090
CMD ["/app/main"] 
ENTRYPOINT ["/app/start.sh"]