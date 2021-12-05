# Terraform Provider Cortex

[![Latest release](https://img.shields.io/github/v/release/inuits/terraform-provider-cortex)](https://github.com/inuits/terraform-provider-cortex/releases)

This provider enables the provisioning of [cortex](https://cortexmetrics.io).

- [Documentation](https://registry.terraform.io/providers/inuits/cortex/latest/docs)

## Manually building the provider

```shell
make build
````

To install into `~/.terraform.d/plugins`:

```shell
make install
```

## Development

## Acceptance tests

Terraform acceptance tests run a local Cortex instance.
```shell
make cortex-up
make testacc
make cortex-down
```

## Ad-hoc testing

For fast feedback loop during development you can run arbitrary Terraform plans using a locally built provider. To run
these plans against dockerized Cortex instance you can use Terraform development overrides.

To enable Terraform logs and configure development overrides:
```shell
make dev.tfrc
source tools/setup-env.sh
```

Build the provider and run some terraform, e.g one of the provided examples:
```shell
make build
terraform init
terraform -chdir=examples apply
```

Remove local Terraform state:
```shell
make clean
```
