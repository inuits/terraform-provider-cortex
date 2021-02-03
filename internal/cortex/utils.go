package cortex

import "gopkg.in/yaml.v3"

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
