.ONESHELL:

TAG=latest

.PHONY: build
build:
	operator-sdk build --go-build-args "-o build/_output/bin/newrelic-alert-manager" fpetkovski/newrelic-alert-manager:$(TAG)

.PHONY: release
release: genapi gendocs build
	docker push fpetkovski/newrelic-alert-manager:$(TAG)

.PHONY: e2etest
e2etest:
	operator-sdk test local ./e2e_tests --up-local --namespace e2e-tests

.PHONY: gendocs
gendocs:
	./hack/docs/gen-crd-api-reference-docs  -template-dir hack/docs/templates -config hack/docs/config.json -api-dir "github.com/fpetkovski/newrelic-alert-manager/pkg/apis/" -out-file docs/README.md

.PHONY: genapi
genapi:
	operator-sdk generate k8s
	operator-sdk generate crds