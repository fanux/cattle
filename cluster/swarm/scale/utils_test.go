package scale

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertKVStringsToMap(t *testing.T) {
	result := convertKVStringsToMap([]string{"HELLO=WORLD", "a=b=c=d", "e"})
	expected := map[string]string{"HELLO": "WORLD", "a": "b=c=d", "e": ""}
	assert.Equal(t, expected, result)
}

func TestConvertMapToKVStrings(t *testing.T) {
	result := convertMapToKVStrings(map[string]string{"HELLO": "WORLD", "a": "b=c=d", "e": ""})
	sort.Strings(result)
	expected := []string{"HELLO=WORLD", "a=b=c=d", "e="}
	assert.Equal(t, expected, result)
}

func TestInaffinityStrings(t *testing.T) {
	res := getInaffinityStrings([]string{"affinity:key!=value", "affinitykey!=value", "affinity:key==value"})
	expected := []string{"key==value"}
	assert.Equal(t, expected, res)
}

func TestConstaintsStrings(t *testing.T) {
	res := getConstaintStrings([]string{"cona:key!=value", "constraint:key!=value", "constraint:key==value", "constraint:key=value"})
	expected := []string{"key!=value", "key==value"}
	assert.Equal(t, expected, res)
}

func TestMatchLabels(t *testing.T) {
	// var res bool
	// var expected bool

	// labels := make(map[string]string)
	// labels["key"] = "value"
	// labels["foo"] = "bar"

	// filters1 := []common.Filter{{"key", "==", "value"}}
	// res = matchAnyLabels(labels, filters1)
	// expected = true
	// assert.Equal(t, expected, res)

	// filters2 := []common.Filter{{"key", "!=", "value"}}
	// res = matchAnyLabels(labels, filters2)
	// expected = false
	// assert.Equal(t, expected, res)
}