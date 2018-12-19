package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/labels"
)

//Rule defines annotation rule
type Rule struct {
	Selector    labels.Set        `yaml:"selector" json:"selector"`
	Annotations map[string]string `yaml:"annotations" json:"annotations"`
}

//LoadRules initialize rules from config source
func LoadRules() (string, bool) {
	if RulesFile == "" {
		return RulesFile, false
	}

	rulesBytes, err := ioutil.ReadFile(RulesFile)
	if err != nil {
		AppLogger.Errorf("Failed to read rules file: %v", err)
	}
	yaml.Unmarshal(rulesBytes, &Rules)

	return RulesFile, true
}
