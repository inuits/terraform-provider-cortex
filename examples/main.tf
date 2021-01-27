terraform {
  required_providers {
    cortex = {
      version = "0.1"
      source  = "inuits/cortex"
    }
  }
}

provider "cortex" {
  address   = "http://cortex"
  api_key   = "eetha9Nahshieghi2xeejih1ethu1aexq"
  tenant_id = "example"
}

resource "cortex_alertmanager_config" "main" {
  content = <<EOT
--
route:
  group_by: ['alertname']
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 1h
  receiver: 'web.hook'
receivers:
- name: 'web.hook'
  webhook_configs:
  - url: 'http://127.0.0.1:5001/'
inhibit_rules:
  - source_match:
      severity: 'critical'
    target_match:
      severity: 'warning'
    equal: ['alertname', 'dev', 'instance']
EOT
}

resource "cortex_alertmanager_ruler_rules" "watchdog" {
  group   = "watchdog"
  content = <<EOT
name: watchdog
rules:
- name: Watchdog
  expr: vector(1)
  for: 10m
EOT
}
