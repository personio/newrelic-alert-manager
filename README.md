## Newrelic alert manager

newrelic-alert-manager is a Kubernetes operator which automates the management of 
Newrelic alert policies and notification channels.

It allows end users of a Kubernetes cluster to define alerting policies and channels as [Kubernetes Custom Resources](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)

### Project status
The project is currently in an alpha version and might not be suitable for production usage.
Additionally, we still cannot offer strong guarantees of API stability. 

### Supported features
The newrelic-alert-manager currently supports the management of the following alerting conditions
* [NRQL alerting conditions](https://docs.newrelic.com/docs/alerts/new-relic-alerts/defining-conditions/create-alert-conditions-nrql-queries)
* [APM alerting conditions](https://docs.newrelic.com/docs/alerts/new-relic-alerts/defining-conditions/create-alert-conditions)

With respect to notification channels, the only supported type is a Slack channel.  

### Deployment
In order to deploy the operator, execute the following steps:

* Clone this repository
* Deploy the custom resource definitions by running
```kubectl apply -f deploy/crds/```
* Add your base64 encoded newrelic admin password to deploy/secret.yaml
* Deploy the operator manifests by running
```kubectl apply -f deploy/```

### Example Usage
Please check the [examples](https://github.com/fpetkovski/newrelic-alert-manager/tree/master/hack/examples) folder to find out how to deploy alert policies together with notification channels.
