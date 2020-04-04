.ONESHELL:

TAG=latest

.PHONY: build
build:
	operator-sdk build --go-build-args "-o build/_output/bin/newrelic-alert-manager" fpetkovski/newrelic-alert-manager:$(TAG)

.PHONY: release
release: build
	docker push fpetkovski/newrelic-alert-manager:$(TAG)