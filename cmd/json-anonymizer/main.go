package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/tleyden/json-anonymizer"
)

func main() {

	// // Anything that starts with an underscore
	// regexpStartsUnderscore, err := regexp.Compile("_(.)*")
	// if err != nil {
	// 	panic(fmt.Sprintf("failed to compile regex.  Err: %v", err))
	// }

	config := json_anonymizer.JsonAnonymizerConfig{
	// SkipFieldsMatchingRegex: []*regexp.Regexp{
	// 	regexpStartsUnderscore,
	// },
	// AnonymizeKeys: true,
	}

	buf := bufio.NewReader(os.Stdin)
	var line []byte

	for {
		raw, isPrefix, err := buf.ReadLine()
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Fatalln(err)
		}
		line = append(line, raw...)
		if isPrefix {
			continue
		}

		jsonAnonymizer := json_anonymizer.NewJsonAnonymizer(config)

		jsonUnmarshalled := map[string]interface{}{}

		if err := json.Unmarshal(line, &jsonUnmarshalled); err != nil {
			panic(fmt.Sprintf("failed to unmarshal json file.  Err: %v", err))
		}

		anonymized, err := jsonAnonymizer.Anonymize(jsonUnmarshalled)
		if err != nil {
			panic(fmt.Sprintf("failed to anonymize json file.  Err: %v", err))

		}

		anonymizedMarshalled, err := json.Marshal(anonymized)
		if err != nil {
			panic(fmt.Sprintf("failed to marshal anonymized json file.  Err: %v", err))
		}

		fmt.Printf("%v\n", string(anonymizedMarshalled))
		line = line[:0]
	}

}
