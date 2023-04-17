

FROM golang:1.19.8-alpine3.16 AS builder

WORKDIR /app

COPY . .

RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-386.tar.gz   | tar xvz

#run

FROM alpine:3.13

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY app.env .
COPY start.sh .
COPY wait.sh .
COPY db/migration ./migration

EXPOSE 8080

CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]
