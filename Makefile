NAME := snooze-otlp
VERSION := 1.0.0
COMMIT := $(shell git rev-parse --short HEAD)

# Default to RHEL bundle. Override it in .make.env
export CA_BUNDLE=/etc/pki/ca-trust/source/anchors/

# Override local dev environment
ifeq ($(shell test -e .make.env && echo yes), yes)
include .make.env
export $(shell sed 's/=.*//' .make.env)
endif

export DOCKER_BUILDKIT := 1
export BUILDAH_FORMAT := docker
export ESC := '\033'
export INFO := "$(ESC)[34m"
export RESET := "$(ESC)[m"

APP_VERSION := v1.0.0
export CHART_VERSION := $(shell grep -oP 'version: \K.*' charts/snooze-otlp/Chart.yaml)

.PHONY: setup develop build release

setup:
	@echo -e "${INFO}Syncing local certificates${RESET}"
	rsync -ahP ${CA_BUNDLE} .ca-bundle/
ifneq ("${LOCAL_REPO}", "")
	@echo -e "${INFO}Authenticate to local repo${RESET}"
	docker login "${LOCAL_REPO}"
endif
ifneq ($(shell test -e .helmfile.yaml && echo yes), yes)
	@echo -e "${INFO}Generate .helmfile.yaml${RESET}"
	NAME=$(NAME) NAMESPACE=$(NAMESPACE) envsubst <.template-helmfile.yaml >.helmfile.yaml
endif

develop:
	@echo -e "${INFO}1) Building docker image${RESET}"
	docker build -t ${LOCAL_REPO}:develop --build-arg COMMIT=$(COMMIT) .
	@echo -e "${INFO}2) Uploading docker image${RESET}"
	docker push ${LOCAL_REPO}:develop
	@echo -e "${INFO}3) Building helm chart${RESET}"
	helm package charts/$(NAME) --app-version develop --version 0.0.0-dev -d .charts/
	@echo -e "${INFO}4) Uploading helm chart${RESET}"
	helm cm-push .charts/$(NAME)-0.0.0-dev.tgz jnx-repo-upload
	@echo -e "${INFO}5) Running helmfile${RESET}"
	helmfile -f .helmfile.yaml sync
