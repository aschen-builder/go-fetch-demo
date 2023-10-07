FROM golang:1.21.0-bullseye

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /usr/bin/fetch

WORKDIR /root

CMD [ "/bin/bash" ]