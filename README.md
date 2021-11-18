# Terraform Provider Cortex

[![Latest release](https://img.shields.io/github/v/release/inuits/terraform-provider-cortex)](https://github.com/inuits/terraform-provider-cortex/releases)

This provider enables the provisioning of [cortex](https://cortexmetrics.io).

- [Documentation](https://registry.terraform.io/providers/inuits/cortex/latest/docs)

## Building the provider

Run the following command to build the provider

```shell
go build -o terraform-provider-cortex
```

## Test sample configuration

First, build and install the provider.

```shell
make install
```

Then, run the following command to initialize the workspace and apply the sample configuration.

```shell
terraform init && terraform apply
```

## Development

For fast feedback loop during development you can use a dockerized Cortex instance and Terraform CLI overrides.
Each step is available as a `make` target.

Available local development targets:

Enable Terraform logs and dev overrides

```shell
make dev.tfrc
source tools/setup-env.sh
```

Run Cortex docker container

```shell
make cortex-up
```

Shut down and cleanup Cortex container

```shell
make cortex-down
```

Remove local Terraform state

```shell
make clean
```