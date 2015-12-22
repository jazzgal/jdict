package main

import (
	"encoding/json"
	"fmt"
	"jdict/lib/jmdict"
	"jdict/lib/kanji"
	"os"
)

// Search & print out result
func searchAndPrint(key string) {
	result := jmdict.Query(key)
	b, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(b))
}

func main() {
	args := os.Args[1:]
	if len(args) > 0 {
		if args[0] == `populate` {
			// Populate db
			if len(args) > 1 {
				fmt.Println("Populate Database ...")
				//jmdict.PopulateData("data/sample.dict")
				jmdict.PopulateData(args[1])
				return

			}
		} else if args[0] == `search` {
			// Query word
			if len(args) > 1 {
				key := args[1]
				searchAndPrint(key)
				return
			}
		} else if args[0] == `populate_kanji` {
			// Query word
			if len(args) > 1 {
				fmt.Println("Populate Kanji Database ...")
				kanji.PopulateData(args[1])
				return
			}
		}

	}
	// Default helper
	fmt.Println(`
Usage: main [ACTION] [OPTIONS]

Actions:
  populate	Parse xml dictionary file & push to database
			Option: xml file path
  search	Query for key
			Option: key string

	`)
}
