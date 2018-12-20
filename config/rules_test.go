package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitRules(t *testing.T) {
	rulesFile := os.Getenv("RULES_FILE")

	os.Setenv("RULES_FILE", "testdata/rules.yaml")
	err := InitRules()
	assert.Nil(t, err)
	assert.Len(t, Rules, 2)
	os.Setenv("RULES_FILE", rulesFile)
}

func TestInitRulesEmpty(t *testing.T) {
	rulesFile := os.Getenv("RULES_FILE")

	os.Setenv("RULES_FILE", "")
	err := InitRules()
	assert.Nil(t, err)
	assert.Len(t, Rules, 0)
	os.Setenv("RULES_FILE", rulesFile)
}

func TestInitRulesError(t *testing.T) {
	rulesFile := os.Getenv("RULES_FILE")

	os.Setenv("RULES_FILE", "testdata/non-existing-file")
	err := InitRules()
	assert.NotNil(t, err)
	assert.Len(t, Rules, 0)
	os.Setenv("RULES_FILE", rulesFile)
}

func TestLoadRules(t *testing.T) {
	rules, err := LoadRules("testdata/rules.yaml")

	assert.Nil(t, err)
	assert.Len(t, rules, 2)
}

func TestLoadRulesEmpty(t *testing.T) {
	rules, err := LoadRules("")

	assert.Nil(t, err)
	assert.Len(t, rules, 0)
}
func TestLoadRulesError(t *testing.T) {
	rules, err := LoadRules("testdata/non-existing-file")

	assert.NotNil(t, err)
	assert.Len(t, rules, 0)
}
