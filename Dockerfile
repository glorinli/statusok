FROM golang:1.17

LABEL org.opencontainers.image.source="https://github.com/vemonet/statusok"

RUN mkdir /app
ADD . /app/
WORKDIR /app

RUN go build
RUN go install

VOLUME /config
ENTRYPOINT ["/app/docker-entrypoint.sh"]
