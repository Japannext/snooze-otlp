APP_VERSION := v1.0.0
export CHART_VERSION := $(shell grep -oP 'version: \K.*' charts/snooze-otlp/Chart.yaml)

.PHONY: build develop release go_build docker_build

go_build:
	mkdir -p build
	CGO_ENABLED=0 GOOS=linux go build -v \
		-ldflags "-X github.com/japannext/snooze-otlp/server.Version=$(APP_VERSION) -w -s" \
		-o build/snooze-otlp-$(APP_VERSION) \
		.
	ln -nsf ./snooze-otlp-$(APP_VERSION) build/snooze-otlp.latest
	cp build/snooze-otlp-$(APP_VERSION) build/snooze-otlp

docker_build:
	docker build .

build: go_build docker_build

develop: go_build
	scripts/develop.sh

release:
	...
