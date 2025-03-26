#!/usr/bin/env bash
set -x

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

pushd "${SCRIPT_DIR}/.."
go mod download

GOMAXPROCS=1 go run "./cmd/orb-rps-server/..." &
server_pid=$!
sleep 5

for i in "drpc" "grpc" "grpcs" "h2c" "http" "https" "http3"; do
  go run "./cmd/orb-rps-client/..." --threads=2 --duration=5 --connections=6 --transport=$i
  sleep 2
done

kill "${server_pid}"
popd
