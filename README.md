## New Relic alert manager

newrelic-alert-manager is a Kubernetes operator which automates the management of 
New Relic dashboards, alert policies and notification channels.

It allows end users of a Kubernetes cluster to define these resources as [Kubernetes Custom Resources](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)

### Project status
The project is currently in an alpha state and might not be suitable for production usage.
Additionally, we still cannot offer strong guarantees of API stability. 

### Supported features
#### Alerts
The newrelic-alert-manager currently supports the management of the following alerting conditions
* [NRQL alerting conditions](https://docs.newrelic.com/docs/alerts/new-relic-alerts/defining-conditions/create-alert-conditions-nrql-queries)
* [APM alerting conditions](https://docs.newrelic.com/docs/alerts/new-relic-alerts/defining-conditions/create-alert-conditions)
* [Infra alerting conditions](https://docs.newrelic.com/docs/infrastructure/new-relic-infrastructure/infrastructure-alert-conditions/rest-api-calls-new-relic-infrastructure-alerts) of type `infra_metric`

If you are unable to create a particular alerting condition due to lack of support by the operator or the New Relic API,
you can try to fall back to defining it as a NRQL alerting condition instead.
One such example is given in the [FAQ](https://github.com/fpetkovski/newrelic-alert-manager#how-do-i-create-an-apm-condition-of-type-web-transaction-percentiles) section. 

#### Notification channels

With respect to notification channels, the currently supported types are Email and Slack channels.  

#### Dashboards
The dashboard API is fully covered by the operator.

### Deployment
In order to deploy the operator, execute the following steps:

* Clone this repository
* Add your base64 encoded New Relic admin password to `deploy/1-secret.yaml`
* Deploy the custom resource definitions by running
```kubectl apply -f deploy/crds/```
* Deploy the operator manifests by running
```kubectl apply -f deploy/```

### Example Usage
Please check the [examples](https://github.com/fpetkovski/newrelic-alert-manager/tree/master/hack/examples) folder to find out how to deploy alert policies together with notification channels.

For more detailed information, the complete API reference can be found [here](https://github.com/fpetkovski/newrelic-alert-manager/tree/master/docs)

### Debugging resources
If you applied an alert policy but it was not created in New Relic, you can check the 
status of the policy using `kubectl describe alertpolicies <policy-name>`. If there was an error while creating the policy, it will be shown in the `Status.reason` field.
Similarly, you can use `kubectl describe` to debug dashboards and notification channels as well.

### FAQ
#### Where can I find a more information on how each alerting condition parameter affects the alert policy?  
The alert condition parameters are best explained by the documentation for the New Relic REST API
Some examples include:
* [apmConditions.alertThreshold.metric](https://docs.newrelic.com/docs/alerts/rest-api-alerts/new-relic-alerts-rest-api/alerts-conditions-api-field-names#metric)
* [apmConditions.alertThreshold.timeFunction](https://docs.newrelic.com/docs/alerts/rest-api-alerts/new-relic-alerts-rest-api/alerts-conditions-api-field-names#terms_time_function)
* [nrqlConditions.valueFunction](https://docs.newrelic.com/docs/alerts/rest-api-alerts/new-relic-alerts-rest-api/alerts-conditions-api-field-names#user_defined_value_function)

You can review the [Alerts conditions API field names](https://docs.newrelic.com/docs/alerts/rest-api-alerts/new-relic-alerts-rest-api/alerts-conditions-api-field-names) page for more information.

#### How do I create an APM condition of type Web transaction percentiles
Unfortunately, it is [not possible](https://docs.newrelic.com/docs/alerts/rest-api-alerts/new-relic-alerts-rest-api/rest-api-calls-new-relic-alerts#excluded) to use New Relic's REST API to create these types of conditions.
However, you can try to define a NRQL alerting condition instead. The query parameter could be defined as follows: 
```
SELECT percentile(totalTime) FROM Transaction WHERE appName = '<your APM application name>'
```
 