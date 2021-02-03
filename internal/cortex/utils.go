package cortex

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gopkg.in/yaml.v3"
)

func formatYAML(input string) (string, error) {
	var rg interface{}
	err := yaml.Unmarshal([]byte(input), &rg)
	if err != nil {
		return "", err
	}
	out, err := yaml.Marshal(rg)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func suppressYAMLDiff(k, old, new string, d *schema.ResourceData) bool {
	olds, err := formatYAML(old)
	if err != nil {
		return false
	}
	news, err := formatYAML(new)
	if err != nil {
		return false
	}
	return olds == news
}
