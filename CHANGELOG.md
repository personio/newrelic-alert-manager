# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).


## [1.5.1] - 2021-01-26
- Fix a bug with updating alert policies when policies span multiple pages  

## [1.5.0] - 2020-11-11
- Add support for text-based widgets using visualization: markdown

## [1.4.0] - 2020-11-07
- Make sinceMinutes in nrqlConditions optional

## [1.3.1] - 2020-11-05
- Add sensible defaults for signal and expiration

## [1.3.0] - 2020-11-05
- Add signal and expiration properties to NRQL alert conditions

## [1.2.1] - 2020-08-05
- Fix a bug with updating notification channels when more than channel records exist in New Relic

## [1.2.0] - 2020-08-05
- Add support notes field in dashboard widgets

## [1.1.0] - 2020-05-29
- Add support for Opsgenie as a notification channel

## [1.0.1] - 2020-05-21
- Eliminates a memory leak caused by the defaults of `http.Client` by setting an explicit connection timeout.

## [1.0.0] - 2020-05-14
- A first stable release of the New Relic alert manager.
- Enables the management of New Relic alerts, notifications channels and dashboards as Kubernetes custom resources.

[1.0.0]: https://github.com/personio/newrelic-alert-manager/releases/tag/v1.0.0
[1.0.1]: https://github.com/personio/newrelic-alert-manager/releases/tag/v1.0.1
[1.1.0]: https://github.com/personio/newrelic-alert-manager/releases/tag/v1.1.0
[1.2.0]: https://github.com/personio/newrelic-alert-manager/releases/tag/v1.2.0
[1.2.1]: https://github.com/personio/newrelic-alert-manager/releases/tag/v1.2.0
