#!/usr/bin/env bash
set -ex

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

for script in $(find $SCRIPT_DIR/../ -name "test.sh"); do
    if [[ "$script" -ef "${BASH_SOURCE[0]}" ]]; then
        continue
    fi

    bash $script
done