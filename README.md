<a id="markdown-AWS Config Filterd" name="AWS Config Filterd"></a>
# AWS Config Filterd
[![GoDoc](https://godoc.org/github.com/asecurityteam/awsconfig-filterd?status.svg)](https://godoc.org/github.com/asecurityteam/awsconfig-filterd)

A Lambda that filters the Config change stream so that only a selection of
relevant resource types continue onto the next stream.

<https://github.com/asecurityteam/awsconfig-filterd>

<!-- TOC -->
- [Overview](#overview)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
    - [Logging](#logging)
    - [Stats](#stats)
- [Supported Resources](#supported-resources)
- [Status](#status)
- [Contributing](#contributing)
    - [Building And Testing](#building-and-testing)
    - [Quality Gates](#quality-gates)
    - [License](#license)
    - [Contributing Agreement](#contributing-agreement)
<!-- /TOC -->

<a id="markdown-overview" name="overview"></a>
## Overview

AWS Config provides a detailed view of the configuration of AWS resources, potentially across multiple AWS accounts, and
can provide a stream of configuration change events via an SNS topic which publishes to SQS. However, much of the data produced by AWS Config may be unnecessary for many use cases.

The awsconfig-filterd service provides a lambda handler which accepts the [configuration item change notification](https://docs.aws.amazon.com/config/latest/developerguide/example-sns-notification.html)
payload, applies filters to remove configuration change events based parameters such as resource type, and returns the
filtered events. The goal is to provide a stream of events for consumption by other services which are only interested
in a subset of the AWS Config data.

<a id="markdown-quick-start" name="quick-start"></a>
## Quick Start

Install docker and docker-compose.

The app can be run locally by running `make run`.

This will run `docker-compose` for the serverfull project
as well as the supplied serverfull-gateway configuration.
The sample configration provided assumes there will be a stats
collector running. To disable this, remove the stats configuration
lines from the server configuration and the serverfull-gateway
configuration.

The app should now be running on port 8080.

`curl -vX POST "http://localhost:8080" -H "Content-Type:application/json" -d @pkg/handlers/v1/testdata/config.valid.json`

<a id="markdown-configuration" name="configuration"></a>
## Configuration

Images of this project are built, and hosted on [DockerHub](https://cloud.docker.com/u/asecurityteam/repository/docker/asecurityteam/awsconfig-filterd). The system is configured using environment variables. The following are all of the configuration options for the system:

```bash
# (bool) Use the Lambda SDK to start the system.
AWSCONFIGFILTERD_LAMBDAMODE="false"
# ([]string) List of enabled filters.
AWSCONFIGFILTERD_FILTER_ENABLED="resourcetype"
# ([]string) List of AWS resource types allowed to pass through.
AWSCONFIGFILTERD_FILTER_RESOURCETYPE_ALLOWED="AWS::EC2::Instance AWS::ElasticLoadBalancing::LoadBalancer AWS::ElasticLoadBalancingV2::LoadBalancer"
# (string)
AWSCONFIGFILTERD_PRODUCER_TYPE="BENTHOS"
# (string) The YAML or JSON text of a Benthos configuration.
AWSCONFIGFILTERD_PRODUCER_BENTHOS_YAML=""
# (string) The URL to POST.
AWSCONFIGFILTERD_PRODUCER_POST_ENDPOINT=""
# (string) The type of HTTP client. Choices are SMART and DEFAULT.
AWSCONFIGFILTERD_PRODUCER_POST_HTTPCLIENT_TYPE="DEFAULT"
# (string) The full OpenAPI specification with transportd extensions.
AWSCONFIGFILTERD_PRODUCER_POST_HTTPCLIENT_SMART_OPENAPI=""
# ([]string) List of enabled lambda decorators
AWSCONFIGFILTERD_DECORATOR_ENABLED="subscription"
# (string) The type of HTTP client. Choices are SMART and DEFAULT.
AWSCONFIGFILTERD_DECORATOR_SUBSCRIPTION_HTTPCLIENT_TYPE="default"
# (string) The full OpenAPI specification with transportd extensions.
AWSCONFIGFILTERD_DECORATOR_SUBSCRIPTION_HTTPCLIENT_SMART_OPENAPI=""
```

For those who do not have access to AWS Lambda, you can run your own configuration by composing this
image with your own custom configuration of serverfull-gateway.

<a id="markdown-logging" name="logging"></a>
### Logging

This project makes use of [logevent](https://github.com/asecurityteam/logevent) which provides structured logging
using Go structs and tags. By default the project will set a logger value in the context for each request. The handler
uses the `LogFn` function defined in `pkg/domain/alias.go` to extract the logger instance from the context.

The built in logger can be configured through the serverfull runtime [configuration](https://github.com/asecurityteam/serverfull#configuration).

<a id="markdown-stats" name="stats"></a>
### Stats

This project uses [xstats](https://github.com/rs/xstats) as its underlying stats library. By default the project will
set a stat client value in the context for each request. The handler uses the `StatFn` function defined in
`pkg/domain/alias.go` to extract the logger instance from the context.

The built in stats client can be configured through the serverfull runtime [configuration](https://github.com/asecurityteam/serverfull#configuration).

Additional resources:

* [serverfull](https://github.com/asecurityteam/serverfull)
* [serverfull-gateway](https://github.com/asecurityteam/serverfull-gateway)

<a id="markdown-supported-resources" name="supported-resources"></a>
## Supported Filters

The current version only supports filtering [configuration item change notification](https://docs.aws.amazon.com/config/latest/developerguide/example-sns-notification.html) events by resource type.

<a id="markdown-status" name="status"></a>
## Status

This project is in incubation which means we are not yet operating this tool in production
and the interfaces are subject to change.

<a id="markdown-contributing" name="contributing"></a>
## Contributing

If you are interested in contributing to the project, feel free to open an issue or PR.

<a id="markdown-building-and-testing" name="building-and-testing"></a>
### Building And Testing

We publish a docker image called [SDCLI](https://github.com/asecurityteam/sdcli) that
bundles all of our build dependencies. It is used by the included Makefile to help make
building and testing a bit easier. The following actions are available through the Makefile:

-   make dep

    Install the project dependencies into a vendor directory

-   make lint

    Run our static analysis suite

-   make test

    Run unit tests and generate a coverage artifact

-   make integration

    Run integration tests and generate a coverage artifact

-   make coverage

    Report the combined coverage for unit and integration tests

-   make build

    Generate a local build of the project (if applicable)

-   make run

    Run a local instance of the project (if applicable)

-   make doc

    Generate the project code documentation and make it viewable
    locally.

<a id="markdown-quality-gates" name="quality-gates"></a>
### Quality Gates

Our build process will run the following checks before going green:

-   make lint
-   make test
-   make integration
-   make coverage (combined result must be 85% or above for the project)

Running these locally, will give early indicators of pass/fail.

<a id="markdown-license" name="license"></a>
### License

This project is licensed under Apache 2.0. See LICENSE.txt for details.

<a id="markdown-contributing-agreement" name="contributing-agreement"></a>
### Contributing Agreement

Atlassian requires signing a contributor's agreement before we can accept a
patch. If you are an individual you can fill out the
[individual CLA](https://na2.docusign.net/Member/PowerFormSigning.aspx?PowerFormId=3f94fbdc-2fbe-46ac-b14c-5d152700ae5d).
If you are contributing on behalf of your company then please fill out the
[corporate CLA](https://na2.docusign.net/Member/PowerFormSigning.aspx?PowerFormId=e1c17c66-ca4d-4aab-a953-2c231af4a20b).
