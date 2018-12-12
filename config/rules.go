package config

import (
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/labels"
)

type Rule struct {
	Selector    labels.Set        `yaml:"selector" json:"selector"`
	Annotations map[string]string `yaml:"annotations" json:"annotations"`
}

//Rules rules
var Rules []Rule

func init() {
	var rulesFile string
	if val, ok := os.LookupEnv("RULES_FILE"); ok {
		rulesFile = val
	} else {
		rulesFile = "config.yaml"
	}
	rulesBytes, err := ioutil.ReadFile(rulesFile)
	if err != nil {
		logrus.Errorf("Failed to read rules file: %v", err)
	}
	yaml.Unmarshal(rulesBytes, &Rules)
}
