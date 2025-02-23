#!/usr/bin/env bash
set -ex; set -o pipefail

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

$SCRIPT_DIR/../benchmarks/event/scripts/test.sh
$SCRIPT_DIR/../benchmarks/rps/scripts/test.sh
$SCRIPT_DIR/../event/simple/scripts/test.sh
$SCRIPT_DIR/../rest/middleware/scripts/test.sh