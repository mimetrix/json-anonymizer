package json_anonymizer

import (
	"testing"
	"log"
	"github.com/stretchr/testify/assert"
)

func TestAnonymize(t *testing.T) {

	config := JsonAnonymizerConfig{
		SkipFieldsMatchingRegex: []string{},
	}
	jsonAnonymizer := NewJsonAnonymizer(config)

	//nestedMap := map[string]interface{} {
	//	"nestedfoo": "nestedbar",
	//}

	testMap := map[string]interface{} {
		"foo": "bar",
		"key": 23423.4,
		"list": []string{ "s1, s2"},
		"nestedMap": map[string]interface{} {
			"nestedfoo": "nestedbar",
		},
		//"nestedlistOfMaps": []map[string]interface{}{
		//	map[string]interface{} {
		//		"nestedlistOfMapsMap1": "nestedbar",
		//	},
		//	map[string]interface{} {
		//		"nestedlistOfMapsMap2": "nestedbar",
		//	},
		//},
		//"nestedlistOStrings": []string{
		//	"s1", "s2",
		//},
	}

	anonymized, err := jsonAnonymizer.Anonymize(testMap)
	if err != nil {
		t.Fatalf("Error anonymizing: %v", err)
	}

	anonymizedMap, ok := anonymized.(map[string]interface{})
	if !ok {
		t.Fatalf("Error type asserting")
	}
	log.Printf("anonymized: %+v", anonymized)

	// we shouldn't have a key with "foo"
	_, hasFooKey := anonymizedMap["foo"]
	assert.False(t, hasFooKey, "Should not have foo key")

	// but we should have key with sha1 hash of foo
	fooValAnon, hasFooKeyAnon := anonymizedMap[anonymizeString("foo")]
	assert.True(t, hasFooKeyAnon, "Should have anonymized foo key")
	assert.Equal(t, fooValAnon, anonymizeString("bar"))

	// we shouldn't have a key with "nestedMap"
	_, hasNestedMapKey := anonymizedMap["nestedMap"]
	assert.False(t, hasNestedMapKey, "Should not have this key")

	// but we should have key with sha1 hash of nestedMap
	nestedMapValAnon, hasNestedMapValAnon := anonymizedMap[anonymizeString("nestedMap")]
	assert.True(t, hasNestedMapValAnon, "Should have anonymized key")

	nestedMapValAnonAsMap := nestedMapValAnon.(map[string]interface{})
	_, hasNestedFooKey := nestedMapValAnonAsMap["nestedfoo"]
	assert.False(t, hasNestedFooKey, "Should not have this key")

	_, hasNestedFooAnonKey := nestedMapValAnonAsMap[anonymizeString("nestedfoo")]
	assert.True(t, hasNestedFooAnonKey, "Should have anonymized key")

	// we shouldn't have a key with "nestedlistOfMaps"
	_, hasNestedListOfMapsKey := anonymizedMap["nestedlistOfMaps"]
	assert.False(t, hasNestedListOfMapsKey, "Should not have this key")




}
