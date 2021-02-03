package cortex

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiffSuppressFunc(t *testing.T) {
	// When the same YAML content is compared, it should return true.

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

	require.True(t, suppressYAMLDiff("", originalYAML, originalYAML, nil))
	require.True(t, suppressYAMLDiff("", originalYAML, equivalentYAML, nil))
	require.False(t, suppressYAMLDiff("", originalYAML, "", nil))
	require.False(t, suppressYAMLDiff("", badYAML, badYAML, nil))
}
