FROM golang:alpine

WORKDIR /build
COPY . /build
RUN go build

CMD ["/build/nanote-server", "/config.yml"]
