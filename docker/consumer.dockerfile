FROM golang:latest

WORKDIR /app

COPY ./consumer .
COPY .env .

RUN make build

CMD cmd/bin/sender
