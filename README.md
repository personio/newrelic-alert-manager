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
* Add your base64 encoded newrelic admin password to `deploy/1-secret.yaml`
* Deploy the custom resource definitions by running
```kubectl apply -f deploy/crds/```
* Deploy the operator manifests by running
```kubectl apply -f deploy/```

### Example Usage
Please check the [examples](https://github.com/fpetkovski/newrelic-alert-manager/tree/master/hack/examples) folder to find out how to deploy alert policies together with notification channels.

### Debugging resources
If you applied an alert policy but it was not created in Newrelic, you can check the 
status of the policy using kubectl describe alertpolicies <policy-name>. If there was an error while creating the policy, it will be shown in the `Status.reason` field.

### FAQ
##### Where can I find a more information on how each alerting condition parameter affects the alert policy?  
The alert condition parameters are best explained by the documentation for the Newrelic REST API
Some examples include:
* [apmConditions.alertThreshold.metric](https://docs.newrelic.com/docs/alerts/rest-api-alerts/new-relic-alerts-rest-api/alerts-conditions-api-field-names#metric)
* [apmConditions.alertThreshold.timeFunction](https://docs.newrelic.com/docs/alerts/rest-api-alerts/new-relic-alerts-rest-api/alerts-conditions-api-field-names#terms_time_function)
* [nrqlConditions.valueFunction](https://docs.newrelic.com/docs/alerts/rest-api-alerts/new-relic-alerts-rest-api/alerts-conditions-api-field-names#user_defined_value_function)

You can review the [Alerts conditions API field names](https://docs.newrelic.com/docs/alerts/rest-api-alerts/new-relic-alerts-rest-api/alerts-conditions-api-field-names) page for more information.

#### How do I create an APM condition of type Web transaction percentiles
Unfortunately, it is [not possible](https://docs.newrelic.com/docs/alerts/rest-api-alerts/new-relic-alerts-rest-api/rest-api-calls-new-relic-alerts#excluded) to use NewRelic's REST API to create these types of conditions.
 