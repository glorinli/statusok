FROM golang:1.11
# FROM golang:1.6.3

ENV GO111MODULE=on

RUN mkdir /app
ADD . /app/
WORKDIR /app

RUN go build
RUN go install

VOLUME /config
ENTRYPOINT ["/app/docker-entrypoint.sh"]
