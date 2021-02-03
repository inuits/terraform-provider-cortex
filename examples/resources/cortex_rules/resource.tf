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
