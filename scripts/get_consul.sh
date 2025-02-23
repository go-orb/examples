#!/usr/bin/env bash
set -e; set -o pipefail

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

# https://gist.github.com/lukechilds/a83e1d7127b78fef38c2914c4ececc3c
function get_latest_gh_release() {
  curl --silent "https://api.github.com/repos/$1/releases/latest" | # Get latest release from GitHub api
    grep '"tag_name":' |                                            # Get tag line
    sed -E 's/.*"([^"]+)".*/\1/'                                    # Pluck JSON value
}

GOOS=$(go env GOOS)
GOARCH=$(go env GOARCH)

VERSION="$(get_latest_gh_release hashicorp/consul)"
VERSION="${VERSION:1}" # Remove the leading v

WORKDIR="$(realpath "${SCRIPT_DIR}/..")/test/bin/${GOOS}_${GOARCH}"

mkdir -p "${WORKDIR}"
pushd "${WORKDIR}" 1>/dev/null 2>&1

if [[ ! -x consul ]]; then
	curl -s -L https://releases.hashicorp.com/consul/${VERSION}/consul_${VERSION}_${GOOS}_${GOARCH}.zip -o consul.zip
	unzip consul.zip 1>/dev/null
	chmod +x consul
	rm -f consul.zip
fi

popd 1>/dev/null 2>&1

echo "${WORKDIR}/consul"
