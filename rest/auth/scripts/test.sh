#!/usr/bin/env bash
set -e; set -o pipefail

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

pushd "${SCRIPT_DIR}/.."
go mod download

go run "./cmd/server/..." 1>/dev/null 2>&1 &
server_pid=$!
sleep 10

go run "./cmd/client/..."

kill "${server_pid}"
popd