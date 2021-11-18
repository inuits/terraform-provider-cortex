terraform {
  required_providers {
    cortex = {
      version = "0.0.4"
      source  = "form3tech-oss/cortex"
    }
  }
}

provider "cortex" {
  address   = "http://127.0.0.1:8080"
  tenant_id = "example"
}

resource "cortex_alertmanager" "main" {
  template_files = {
    default_template = <<EOT
    {{ define "__alertmanager" }}AlertManager{{ end }}
    {{ define "__alertmanagerURL" }}{{ .ExternalURL }}/#/alerts?receiver={{ .Receiver | urlquery }}{{ end }}
EOT
  }
  alertmanager_config = <<EOT
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

resource "cortex_rules" "watchdog" {
  namespace = "default"
  content   = <<EOT
name: watchdog
rules:
- alert: Watchdog
  expr: vector(1)
  for: 20m
EOT
}

resource "cortex_rules" "alerting_rules" {
  namespace = "default"
  for_each  = fileset(path.module, "*_alerts.yml")
  content   = templatefile("${path.module}/${each.key}", {})
}
