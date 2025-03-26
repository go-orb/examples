#!/usr/bin/env bash
set -x

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

pushd "${SCRIPT_DIR}/.."
go mod download

go run "./cmd/foobar/..." server &
server_pid=$!

kill "${server_pid}"
popd