# Terraform Provider Cortex

This provider enables the provisioning of [cortex](https://cortexmetrics.io).


** Status: WIP **

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
