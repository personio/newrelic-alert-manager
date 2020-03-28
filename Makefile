.ONESHELL:

TAG=latest

.PHONY: build
build:
	operator-sdk build fpetkovski/newrelic-alert-manager:$(TAG)

.PHONY: release
release: build
	docker push fpetkovski/newrelic-alert-manager:$(TAG)