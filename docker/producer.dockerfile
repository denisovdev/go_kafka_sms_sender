FROM golang:latest

WORKDIR /app

COPY ./producer .
COPY .env .

RUN make build

EXPOSE 8080

CMD .cmd/bin/messager
