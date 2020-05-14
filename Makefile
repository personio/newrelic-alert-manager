.ONESHELL:

TAG=latest

.PHONY: build
build:
	operator-sdk build --go-build-args "-o build/_output/bin/newrelic-alert-manager" personio/newrelic-alert-manager:$(TAG)

.PHONY: release
release: genapi unittest e2etest gendocs build
	docker push personio/newrelic-alert-manager:$(TAG)

.PHONY: e2etest
e2etest:
	kubectl create ns e2e-tests
	operator-sdk test local ./e2e_tests --up-local --namespace e2e-tests

.PHONY: e2etest-clean
e2etest-clean:
	kubectl delete ns e2e-tests

.PHONY: unittest
unittest:
	go test ./pkg/...

.PHONY: gendocs
gendocs:
	./hack/docs/gen-crd-api-reference-docs  -template-dir hack/docs/templates -config hack/docs/config.json -api-dir "github.com/personio/newrelic-alert-manager/pkg/apis/" -out-file docs/README.md

.PHONY: genapi
genapi:
	operator-sdk generate k8s
	operator-sdk generate crds