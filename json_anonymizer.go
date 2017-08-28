package json_anonymizer

import (
	"github.com/dustin/gojson"
	"log"
	"fmt"
	"crypto/sha1"
	"regexp"
)

type JsonAnonymizer struct {
	Config JsonAnonymizerConfig
}

type JsonAnonymizerConfig struct {
	SkipFieldsMatchingRegex []*regexp.Regexp
	AnonymizeKeys bool
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

		// Rather than add to the map directly while iterating, which can cause "double-anonymizing" certain
		// entries due to the iterator visiting the value and *then* visiting the anonymized value,
		// add any new vals to a separate map.
		newMapVals := map[string]interface{}{}

		for key, val := range v {
			log.Printf("key: %v, val: %v", key, val)

			if ja.ShouldSkip(key) {
				continue
			}

			anonymizedVal, err := ja.Anonymize(val)
			if err != nil {
				return nil, err
			}
			delete(v, key)
			newKey := key
			if (ja.Config.AnonymizeKeys) {
				newKey = anonymizeString(key)
			}

			// Add it to newMapvals to avoid double-anonymizing
			newMapVals[newKey] = anonymizedVal
		}

		// Add all the values from newMapVals into target
		for newKey, newVal := range newMapVals {
			v[newKey] = newVal
		}

		return v, nil
	case []interface{}:
		newSlice := []interface{}{}
		for i, val := range v {
			log.Printf("array index: %v, val: %v", i, val)
			anonymizedVal, err := ja.Anonymize(val)
			if err != nil {
				return nil, err
			}
			newSlice = append(newSlice, anonymizedVal)
		}
		return newSlice, nil
	case float64:
		log.Printf("float64: %v", v)
		return anonymizeFloat64(v), nil
	case string:
		log.Printf("string: %v", v)
		return anonymizeString(v), nil
	case bool:
	case nil:
		// ignore it
	default:
		return nil, err
	}

	return anonymized, nil

}

func (ja JsonAnonymizer) ShouldSkip(key string) bool {
	for _, skipRegexp := range ja.Config.SkipFieldsMatchingRegex {
		if skipRegexp.MatchString(key) {
			return true
		}
	}
	return false
}

//func (ja JsonAnonymizer) AnonymizeMap(input map[string]interface{}) (err error) {
//
//	for key, val := range input {
//
//		switch v := val.(type) {
//		case map[string]interface{}:
//			delete(input, key)
//			// recursively anonymize map
//			if errAnonymizeMap := ja.AnonymizeMap(v); errAnonymizeMap != nil {
//				return errAnonymizeMap
//			}
//			input[anonymizeString(key)] = v
//		case []interface{}:
//			// return fmt.Errorf("Cannot handle slice/list values")
//			for index, listItem := range v {
//
//			}
//		case float64:
//			delete(input, key)
//			input[anonymizeString(key)] = anonymizeFloat64(v)
//		case string:
//			delete(input, key)
//			input[anonymizeString(key)] = anonymizeString(v)
//		case bool:
//		case nil:
//			// ignore it
//		default:
//			return fmt.Errorf("Unknown primitive type: %T.  Val: %v for Key: %v", v, val, key)
//		}
//
//
//	}
//	return nil
//}

func anonymizeString(s string) string {

	// calculate raw sha1 hash
	shaBytes := sha1.Sum([]byte(s))

	// return hex output
	hex := fmt.Sprintf("%x", shaBytes)

	log.Printf("anonymize: %v -> %v", s, hex)
	if (s == "39824c8d5eea6ceca7681985e6b675ef5dd0981e") {
		log.Printf("looks like double anon")
	}

	return hex
}

func anonymizeFloat64(f float64) float64 {
	return 0
}
