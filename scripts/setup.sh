#!/usr/bin/env bash

BASEDIR="$(realpath "${BASH_SOURCE%/*}/..")"

if [[ -f "${BASEDIR}/.env" ]]; then
  # shellcheck source=/dev/null
  source "${BASEDIR}/.env"
fi

CA_PATHS=(
  "/etc/pki/ca-trust/source/anchors" # RHEL
  "/usr/local/share/ca-certificates" # Ubuntu / Alpine
  "/usr/share/pki/trust/anchors" # OpenSUSE
  "/etc/pki/trust/anchors" # OpenSUSE
)

function sync_local_ca () {
  for path in "${CA_PATHS[@]}"; do
    if [[ -d "${path}" ]] && find "${path}" -mindepth 1 -maxdepth 1 | read -r; then
      rsync -ahP "${path}/" "${BASEDIR}/.ca-bundle/"
      return
    fi
  done
  (>&2 echo "Could not find any standard system CA directory")
}
