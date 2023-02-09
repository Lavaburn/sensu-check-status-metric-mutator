[![Sensu Bonsai Asset](https://img.shields.io/badge/Bonsai-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/Lavaburn/sensu-check-status-metric-mutator)
![Go Test](https://github.com/Lavaburn/sensu-check-status-metric-mutator/workflows/Go%20Test/badge.svg)
![goreleaser](https://github.com/Lavaburn/sensu-check-status-metric-mutator/workflows/goreleaser/badge.svg)

# Sensu Check Status Metric Mutator

## Table of Contents
- [Overview](#overview)
- [Usage examples](#usage-examples)
- [Configuration](#configuration)
  - [Asset registration](#asset-registration)
  - [Mutator definition](#mutator-definition)
- [Installation from source](#installation-from-source)
- [Contributing](#contributing)

## Overview

The Sensu Check Status Metric Mutator is a [Sensu Mutator][2] that surfaces the
[exit status of a Sensu Check][6] in [Sensu metric format][7] to be handled
by an output metric handler (such as Influxdb).

## Usage examples

```
Sensu Check Status Metric Mutator

Usage:
  sensu-check-status-metric-mutator [flags]
  sensu-check-status-metric-mutator [command]

Available Commands:
  help        Help about any command
  version     Print the version number of this plugin

Flags:
  -h, --help                          help for sensu-check-status-metric-mutator
  -t, --metric-name-template string   Template for naming the metric point for the check status (default "{{.Check.Name}}.status")

Use "sensu-check-status-metric-mutator [command] --help" for more information about a command.
```

Using this mutator will create a metric point similar to the one below that contains the 
exit status of the check command ran for the event.  The only argument for the mutator is
the template to use for naming the metric point.

```json
  "metrics": {
    "handlers": [
      "influxdb"
    ],
    "points": [
      {
        "name": "check_cpu.status",
        "value": 0,
        "timestamp": 1590369892,
        "tags": [
          {
            "name": "entity",
            "value": "entity1"
          },
          {
            "name": "check",
            "value": "check_cpu"
          },
          {
            "name": "state",
            "value": "passing"
          },
          {
            "name": "occurrences",
            "value": "2958"
          },
          {
            "name": "occurrences_watermark",
            "value": "2958"

      }
    ]
```

## Configuration

### Asset registration

[Sensu Assets][5] are the best way to make use of this plugin. If you're not using an asset, please
consider doing so! If you're using sensuctl 5.13 with Sensu Backend 5.13 or later, you can use the
following command to add the asset:

```
sensuctl asset add Lavaburn/sensu-check-status-metric-mutator
```

If you're using an earlier version of sensuctl, you can find the asset on the [Bonsai Asset Index][3]

### Mutator definition

```yml
---
type: Mutator
api_version: core/v2
metadata:
  name: sensu-check-status-metric-mutator
  namespace: default
spec:
  command: sensu-check-status-metric-mutator
  runtime_assets:
  - Lavaburn/sensu-check-status-metric-mutator
```

This mutator would then be referenced by your metrics handler definition.

```yml
---
type:
api_version; core/v2
metadata:
  name: influxdb
  namespace: default
spec:
  type: pipe
  command: sensu-influxdb-handler -d $INFLUXDB_DB
  timeout: 10
  filters:
  - has_metrics
  mutator: sensu-check-status-metric-mutator
  runtime_assets:
  - sensu/sensu-influxdb-handler
  secrets:
  - name: INFLUXDB_ADDR
    secret: influxdb_addr
  - name: INFLUXDB_DB
    secret: influxdb_db
  - name: INFLUXDB_USER
    secret: influxdb_user
  - name: INFLUXDB_PASSWORD
    secret: influxdb_password
```

## Installation from source

The preferred way of installing and deploying this plugin is to use it as an Asset. If you would
like to compile and install the plugin from source or contribute to it, download the latest version
or create an executable script from this source.

From the local path of the sensu-check-status-metric-mutator repository:

```
go build
```

## 0.4.0 - Forked from nixwiz

This version adds many flags to customize the metrics output.

I wanted to use this for outputting all check results to Prometheus Push Gateway.

* The -occurrences- as a tag (Prometheus label) creates a new series in Prometheus on every push/scrape. It makes a lot more sense to add this data as a separate metric.
* Push Gateway groups metrics over specific labels, hence I wanted the option to have specific tag (label) names.
* For this project I also forked sensu-prometheus-pushgateway-handler, as I only wanted a single metric with labels for entity and check name.
* To enable backwards-compatibility, I used arguments which can be set to "true" or "false" instead of using bools (flags).

```
Flags:
  -h, --help                                      help for sensu-check-status-metric-mutator

  -t, --metric-name-template string               Template for naming the metric point for the check status (default "{{.Check.Name}}.status")

  -c, --tag-name-check string                     The tag name that contains the check name (default "check")
  -e, --tag-name-entity string                    The tag name that contains the entity name (default "entity")
  -s, --tag-name-state string                     The tag name that contains the state name (default "state")

  -m, --show-occurrences-as-metric string         Whether to add a metric for Occurrences (default "false")
  -u, --metric-name-template-occurrences string   Template for naming the metric point for the check status (default "{{.Check.Name}}.occurrences")

  -r, --show-occurrences-as-tag string            Whether to add a tag for Occurrences (default "true")
  -o, --tag-name-occurrences string               The tag name that contains the occurrences name (default "occurrences")

  -w, --enable-watermark string                   The tag name that contains the occurrences name (default "true")
```

## Contributing

For more information about contributing to this plugin, see [Contributing][1].

[1]: https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md
[2]: https://docs.sensu.io/sensu-go/latest/reference/mutators/
[3]: https://bonsai.sensu.io/assets/Lavaburn/sensu-check-status-metric-mutator
[9]: https://github.com/sensu-community/sensu-plugin-tool
[5]: https://docs.sensu.io/sensu-go/latest/reference/assets/
[6]: https://docs.sensu.io/sensu-go/latest/reference/checks/#check-result-specification
[7]: https://docs.sensu.io/sensu-go/latest/reference/events/#metrics
