#!/bin/sh
set -e
sleep 5 # Workaround to wait untill InfluxDb will start

/go/bin/statusok --config /config/config.json
# main --config /config/config.json