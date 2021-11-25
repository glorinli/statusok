FROM golang:1.17
# FROM golang:1.11
# FROM golang:1.6.3

# ENV GO111MODULE=on

RUN mkdir /app
ADD . /app/
WORKDIR /app

# RUN go get github.com/influxdata/influxdb1-client
# RUN go install github.com/influxdata/influxdb1-client
# RUN go install github.com/influxdata/influxdb1-client@latest
# RUN go get gopkg.in/yaml.v2

RUN go build
RUN go install

VOLUME /config
ENTRYPOINT ["/app/docker-entrypoint.sh"]
