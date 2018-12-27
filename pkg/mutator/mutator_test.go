package mutator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatePatchFromAnnotations(t *testing.T) {
	podAnnotations := map[string]string{
		"hello": "world",
	}
	rulesAnnotations := map[string]string{
		"log": "enabled",
	}
	patch := createPatchFromAnnotations(podAnnotations, rulesAnnotations)
	assert.Equal(t, "replace", patch.Op)
	assert.Equal(t, "/metadata/annotations", patch.Path)
	assert.IsType(t, podAnnotations, patch.Value)

	valuesAsMap := patch.Value.(map[string]string)
	assert.Len(t, valuesAsMap, 2)
	assert.Equal(t, "world", valuesAsMap["hello"])
	assert.Equal(t, "enabled", valuesAsMap["log"])
}

func TestCreatePatchFromNilAnnotations(t *testing.T) {
	rulesAnnotations := map[string]string{
		"log": "enabled",
	}
	patch := createPatchFromAnnotations(nil, rulesAnnotations)
	assert.Equal(t, "add", patch.Op)
	assert.Equal(t, "/metadata/annotations", patch.Path)
	assert.IsType(t, make(map[string]string), patch.Value)

	valuesAsMap := patch.Value.(map[string]string)
	assert.Len(t, valuesAsMap, 1)
	assert.Equal(t, "enabled", valuesAsMap["log"])
}
