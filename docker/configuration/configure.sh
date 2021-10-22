#!/usr/bin/env bash

set -ex
## U CAN CHANGE IT IN DOCKER ENV
CONSUL_HOST=${CONSUL_HOST:=consul}
CONSUL_PORT=${CONSUL_PORT:=8500}

wait-for-it.sh --host=${CONSUL_HOST} --port=${CONSUL_PORT} --timeout=300 -q
