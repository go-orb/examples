#!/usr/bin/env bash
set -x

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

pushd "${SCRIPT_DIR}/.."
go mod download

go run "./cmd/server/..." --config config/grpc.yaml 1>/dev/null 2>&1 &
server_pid=$!
sleep 10

go run "./cmd/client/..." --config config/grpc.yaml

kill "${server_pid}"
popd