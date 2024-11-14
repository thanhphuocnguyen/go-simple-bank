# build stage
FROM golang:1.23.3-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY app.env .
COPY db/migrations ./migrations
COPY start.sh .
RUN chmod +x /app/start.sh

EXPOSE 8082
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]