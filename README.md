[![Tests](https://github.com/netascode/terraform-provider-secureaccess/actions/workflows/test.yml/badge.svg)](https://github.com/netascode/terraform-provider-secureaccess/actions/workflows/test.yml)

# Terraform Provider Cisco Secure Access

The Secure Access provider provides resources to interact with a Cisco Secure Access (CSA) over the REST API.

Documentation: <https://registry.terraform.io/providers/netascode/secureaccess/latest>

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.25

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

This Terraform Provider is available to install automatically via `terraform init`. If you're building the provider, follow the instructions to
[install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin)
After placing it into your plugins directory, run `terraform init` to initialize it.

Additional documentation, including available resources and their arguments/attributes can be found on the [Terraform documentation website](https://registry.terraform.io/providers/netascode/secureaccess/latest/docs).

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

## Acceptance tests

Note: Acceptance tests create real resources. You'd need an Secure Access instance with an administrative user on the default global domain.

Depending on whether the environment is full or partial, the suite of Acceptance tests will
execute all/partial tests:

```shell
make testacc
```
