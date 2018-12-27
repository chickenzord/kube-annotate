package config

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/labels"
)

//Rule defines annotation rule
type Rule struct {
	Selector    labels.Set        `yaml:"selector" json:"selector"`
	Annotations map[string]string `yaml:"annotations" json:"annotations"`
}

//LoadRules load rules from config source
func LoadRules(path string) ([]Rule, error) {
	rules := make([]Rule, 0)
	if path == "" {
		return rules, nil
	}

	rulesBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	yaml.Unmarshal(rulesBytes, &rules)

	return rules, nil
}

//InitRules initialize rules from config source
func InitRules() error {
	Rules = make([]Rule, 0)
	rulesFile := os.Getenv("RULES_FILE")

	if len(rulesFile) == 0 {
		AppLogger.Warn("no rules file set")
		return nil
	}

	rules, err := LoadRules(rulesFile)
	if err != nil {
		return err
	}

	AppLogger.Infof("loaded %d rule(s) from %s", len(rules), rulesFile)
	Rules = rules
	return nil
}
