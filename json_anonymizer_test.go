package json_anonymizer

import (
	"encoding/json"
	"testing"

	"log"

	"regexp"

	"github.com/stretchr/testify/assert"
)

func TestAnonymize(t *testing.T) {

	// Anything that starts with an underscore
	regexpStartsUnderscore, err := regexp.Compile("_(.)*")
	if err != nil {
		t.Fatalf("Error compiling regex: %v", err)
	}

	config := JsonAnonymizerConfig{
		SkipFieldsMatchingRegex: []*regexp.Regexp{
			regexpStartsUnderscore,
		},
		AnonymizeKeys: true,
	}
	jsonAnonymizer := NewJsonAnonymizer(config)

	testMap := map[string]interface{}{
		"_id":      "doc1",
		"_rev":     "1-3lkj",
		"_deleted": true,
		"foo":      "bar",
		"key":      23423.4,
		"list":     []string{"s1", "s2"},
		"nestedMap": map[string]interface{}{
			"nestedfoo": "nestedbar",
		},
		"nestedlistOfMaps": []map[string]interface{}{
			map[string]interface{}{
				"nestedlistOfMapsMap1": "nestedbar",
			},
			map[string]interface{}{
				"nestedlistOfMapsMap2": "nestedbar",
			},
		},
		"nestedlistOStrings": []string{
			"s1", "s2",
		},
	}

	anonymized, err := jsonAnonymizer.Anonymize(testMap)
	if err != nil {
		t.Fatalf("Error anonymizing: %v", err)
	}

	anonymizedStr, err := json.MarshalIndent(anonymized, "", "    ")
	if err != nil {
		t.Fatalf("Error marshalling: %v", err)
	}
	log.Printf("anonymized: %+v", string(anonymizedStr))

	anonymizedMap, ok := anonymized.(map[string]interface{})
	if !ok {
		t.Fatalf("Error type asserting")
	}

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
	anonKey := anonymizeString("nestedMap")
	nestedMapValAnon, hasNestedMapValAnon := anonymizedMap[anonKey]
	assert.True(t, hasNestedMapValAnon, "Should have anonymized key: %v", anonKey)

	nestedMapValAnonAsMap := nestedMapValAnon.(map[string]interface{})
	_, hasNestedFooKey := nestedMapValAnonAsMap["nestedfoo"]
	assert.False(t, hasNestedFooKey, "Should not have this key")

	_, hasNestedFooAnonKey := nestedMapValAnonAsMap[anonymizeString("nestedfoo")]
	assert.True(t, hasNestedFooAnonKey, "Should have anonymized key")

	// we shouldn't have a key with "nestedlistOfMaps"
	_, hasNestedListOfMapsKey := anonymizedMap["nestedlistOfMaps"]
	assert.False(t, hasNestedListOfMapsKey, "Should not have this key")

	// Make sure it has fields that should be skipped by anonymizer
	_, hasId := anonymizedMap["_id"]
	assert.True(t, hasId, "Should have this key")
	_, hasRev := anonymizedMap["_rev"]
	assert.True(t, hasRev, "Should have this key")
	_, hasDeleted := anonymizedMap["_deleted"]
	assert.True(t, hasDeleted, "Should have this key")


}
