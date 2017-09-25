package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/tleyden/json-anonymizer"
)

func main() {

	// Anything that starts with an underscore
	regexpStartsUnderscore, err := regexp.Compile("_(.)*")
	if err != nil {
		panic(fmt.Sprintf("failed to compile regex.  Err: %v", err))
	}

	config := json_anonymizer.JsonAnonymizerConfig{
		SkipFieldsMatchingRegex: []*regexp.Regexp{
			regexpStartsUnderscore,
		},
		AnonymizeKeys: true,
	}
	jsonAnonymizer := json_anonymizer.NewJsonAnonymizer(config)

	filename := os.Args[1]

	file, err := os.Open(filename)
	if err != nil {
		panic(fmt.Sprintf("failed to open file.  Err: %v", err))
	}

	jsonUnmarshalled := map[string]interface{}{}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(fmt.Sprintf("failed to read file.  Err: %v", err))
	}

	if err := json.Unmarshal(fileBytes, &jsonUnmarshalled); err != nil {
		panic(fmt.Sprintf("failed to unmarshal json file.  Err: %v", err))
	}

	anonymized, err := jsonAnonymizer.Anonymize(jsonUnmarshalled)
	if err != nil {
		panic(fmt.Sprintf("failed to anonymize json file.  Err: %v", err))

	}

	anonymizedMarshalled, err := json.MarshalIndent(anonymized, "", "    ")
	if err != nil {
		panic(fmt.Sprintf("failed to marshal anonymized json file.  Err: %v", err))
	}

	fmt.Printf("%v", string(anonymizedMarshalled))

}
