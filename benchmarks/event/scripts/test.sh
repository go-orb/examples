#!/usr/bin/env bash
set -ex; set -o pipefail

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

NATS_SERVER=$($SCRIPT_DIR/../../../scripts/get_nats.sh)

$NATS_SERVER -js 1>/dev/null 2>&1 &
nats_pid=$!
sleep 10

pushd "${SCRIPT_DIR}/.."
go mod download

go run "./cmd/handler/..." &
handler_pid=$!
sleep 10

go run "./cmd/request/..." --threads=2 --duration=5 --connections=6

kill "${handler_pid}" "${nats_pid}"
popd