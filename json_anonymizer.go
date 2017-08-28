package json_anonymizer

import (
	"github.com/dustin/gojson"
	"log"
	"fmt"
	"crypto/sha1"
)

type JsonAnonymizer struct {
	Config JsonAnonymizerConfig
}

type JsonAnonymizerConfig struct {
	SkipFieldsMatchingRegex []string
}

func NewJsonAnonymizer(config JsonAnonymizerConfig) *JsonAnonymizer {
	return &JsonAnonymizer{
		Config: config,
	}
}

func (ja JsonAnonymizer) Anonymize(input interface{}) (anonymized interface{}, err error) {

	// Copy the input to the output by marshalling to json and unmarshaling
	inputMarshalled, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(inputMarshalled, &anonymized); err != nil {
		return nil, err
	}

	switch v := anonymized.(type) {
	case map[string]interface{}:
		log.Printf("%v is a map", v)
		if err := ja.AnonymizeMap(v); err != nil {
			return nil, err
		}
	case []interface{}:
		log.Printf("%v is a slice", v)
		return nil, fmt.Errorf("Top level slices/lists not supported yet")
	default:
		return nil, err
	}

	return anonymized, nil

}

func (ja JsonAnonymizer) AnonymizeMap(input map[string]interface{}) (err error) {

	for key, val := range input {

		switch v := val.(type) {
		case map[string]interface{}:
			delete(input, key)
			// recursively anonymize map
			if errAnonymizeMap := ja.AnonymizeMap(v); errAnonymizeMap != nil {
				return errAnonymizeMap
			}
			input[anonymizeString(key)] = v
		case []interface{}:
			// return fmt.Errorf("Cannot handle slice/list values")
			for index, listItem := range v {

			}
		case float64:
			delete(input, key)
			input[anonymizeString(key)] = anonymizeFloat64(v)
		case string:
			delete(input, key)
			input[anonymizeString(key)] = anonymizeString(v)
		case bool:
		case nil:
			// ignore it
		default:
			return fmt.Errorf("Unknown primitive type: %T.  Val: %v for Key: %v", v, val, key)
		}


	}
	return nil
}

func anonymizeString(s string) string {

	// calculate raw sha1 hash
	shaBytes := sha1.Sum([]byte(s))

	// return hex output
	return fmt.Sprintf("%x", shaBytes)

}

func anonymizeFloat64(f float64) float64 {
	return 0
}
