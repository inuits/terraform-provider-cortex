package cortex

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSuppressYAMLDiff(t *testing.T) {
	originalYAML := `route:
  group_by: ['alertname']
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 1h
  receiver: 'web.hook'
`
	equivalentYAML := `route:
  group_by:          ['alertname'    ]
  group_wait:      30s
  group_interval: "5m"
  repeat_interval:      1h
  receiver: web.hook
`
	badYAML := `route:
group_wait:      30s
- group_interval: "5m
`
	tests := []struct {
		name             string
		oldValue         string
		newValue         string
		expectedSuppress bool
	}{
		{"original vs original", originalYAML, originalYAML, true},
		{"original vs equivalent", originalYAML, equivalentYAML, true},
		{"original vs empty", originalYAML, "", false},
		{"bad vs bad", badYAML, badYAML, false},
		{"original vs bad", originalYAML, badYAML, false},
		{"boolean vs string equivalent", `some_bool: true`, `some_bool: "true"`, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expectedSuppress, suppressYAMLDiff("", test.oldValue, test.newValue, nil))
		})
	}
}

func TestSuppressRuleGroupDiff(t *testing.T) {
	ruleGroupYAMLWithBooleanAsBoolean := `name: test
rules:
- alert: test
  expr: 'test > 0'
  labels:
    test_1: true
    severity: critical
    team: platform
  annotations:
    title: 'some title'
    summary: 'some summary'
    description: |-
      Some description
      Link https://github.com
`
	ruleGroupYAMLWithBooleanAsString := `name: test
rules:
- alert: test
  expr: 'test > 0'
  labels:
    test_1: "true"
    severity: critical
    team: platform
  annotations:
    title: 'some title'
    summary: 'some summary'
    description: |-
      Some description
      Link https://github.com
`
	tests := []struct {
		name             string
		oldValue         string
		newValue         string
		expectedSuppress bool
	}{
		{"boolean vs string equivalent", `some_bool: true`, `some_bool: "true"`, true},
		{"RuleGroup with boolean as boolean vs RuleGroup with boolean as string", ruleGroupYAMLWithBooleanAsBoolean, ruleGroupYAMLWithBooleanAsString, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expectedSuppress, suppressRuleGroupDiff("", test.oldValue, test.newValue, nil))
		})
	}
}
