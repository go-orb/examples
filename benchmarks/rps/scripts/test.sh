#!/usr/bin/env bash
set -ex; set -o pipefail

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

CONSUL_SERVER=$($SCRIPT_DIR/../../../scripts/get_consul.sh)

$CONSUL_SERVER agent -bind 127.0.0.1 -data-dir=/tmp/consul -server -bootstrap-expect=1 1>/dev/null 2>&1 &
consul_pid=$!
sleep 5

pushd "${SCRIPT_DIR}/.."
go mod download

GOMAXPROCS=1 go run "./cmd/orb-rps-server/..." --registry=consul &
server_pid=$!
sleep 20

for i in "drpc" "grpc" "h2c" "http" "https" "http3"; do
  go run "./cmd/orb-rps-client/..." --registry=consul --threads=2 --duration=5 --connections=6 --transport=$i
  sleep 2
done

kill "${server_pid}" "${consul_pid}"
popd
