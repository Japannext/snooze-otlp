#!/usr/bin/env bash

set -e

CHART_NAME="snooze-otlp"

if [[ -f .env ]]; then
  source .env
fi

function require_env {
  for var in "${@}"; do
    if [[ -z "${!var}" ]]; then
      (>&2 echo "Missing environment variable '${var}'")
      (>&2 echo "Define it in '.env' or as an environment variable")
      exit 2
    else
      (>&2 echo "${var}=${!var}")
    fi
  done
}

function upload_docker_image {
  image="${1}"
  docker build . -t "${image}"
  docker push "${image}"
}

function upload_helm_chart {
  chart_name="${1}"
  chart_version="${2}"
  chart_repo="${3}"
  mkdir -p build/helm
  helm package "charts/${chart_name}" --version "${chart_version}" --app-version develop -d build/helm
  helm cm-push "build/helm/${chart_name}-${chart_version}.tgz" "${chart_repo}"
}

require_env DOCKER_REPO IMAGE
upload_docker_image "${DOCKER_REPO}/${IMAGE}:develop"

require_env CHART_NAME CHART_VERSION CHART_UPLOAD_REPO
upload_helm_chart "${CHART_NAME}" "${CHART_VERSION}-develop" "${CHART_UPLOAD_REPO}"

require_env CHART_REPO
helmfile -f ".helmfile.yaml" sync
