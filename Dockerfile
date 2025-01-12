FROM golang:1.23.2-alpine

RUN go version

ENV GOPATH=/

COPY . .

RUN go mod download

RUN go build -o app ./cmd

CMD ["./app"]