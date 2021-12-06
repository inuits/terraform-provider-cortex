module github.com/inuits/terraform-provider-cortex

go 1.16

replace k8s.io/client-go => k8s.io/client-go v0.21.0

replace k8s.io/api => k8s.io/api v0.21.0

replace github.com/hashicorp/consul => github.com/hashicorp/consul v1.8.1

replace github.com/bradfitz/gomemcache => github.com/themihai/gomemcache v0.0.0-20180902122335-24332e2d58ab

require (
	github.com/agl/ed25519 v0.0.0-20170116200512-5312a6153412 // indirect
	github.com/alcortesm/tgz v0.0.0-20161220082320-9c5fe88206d7 // indirect
	github.com/grafana/cortex-tools v0.10.4
	github.com/hashicorp/terraform-exec v0.15.0 // indirect
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.9.0
	github.com/keybase/go-crypto v0.0.0-20161004153544-93f5b35093ba // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/prometheus/prometheus v1.8.2-0.20210510213326-e313ffa8abf6
	github.com/stretchr/testify v1.7.0
	github.com/zclconf/go-cty v1.10.0 // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)
