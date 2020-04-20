.ONESHELL:

TAG=latest

.PHONY: build
build:
	operator-sdk build --go-build-args "-o build/_output/bin/newrelic-alert-manager" fpetkovski/newrelic-alert-manager:$(TAG)

.PHONY: release
release: genapi unittest e2etest gendocs build
	docker push fpetkovski/newrelic-alert-manager:$(TAG)

.PHONY: e2etest
e2etest:
	operator-sdk test local ./e2e_tests --up-local --namespace e2e-tests

.PHONY: e2etest-clean
e2etest-clean:
	kubectl delete --all alertpolicies.alerts.newrelic.io -n e2e-tests
	kubectl delete --all dashboards.dashboards.newrelic.io -n e2e-tests
	kubectl delete --all emailnotificationchannels.alerts.newrelic.io -n e2e-tests
	kubectl delete --all slacknotificationchannels.alerts.newrelic.io -n e2e-tests

.PHONY: unittest
unittest:
	go test ./pkg/...

.PHONY: gendocs
gendocs:
	./hack/docs/gen-crd-api-reference-docs  -template-dir hack/docs/templates -config hack/docs/config.json -api-dir "github.com/fpetkovski/newrelic-alert-manager/pkg/apis/" -out-file docs/README.md

.PHONY: genapi
genapi:
	operator-sdk generate k8s
	operator-sdk generate crds