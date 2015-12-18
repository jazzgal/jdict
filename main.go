package main

import (
	"encoding/json"
	"fmt"
	"jdict/lib/jmdict"
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
			fmt.Println("Populate Database ...")
			jmdict.PopulateData("data/sample.dict")
			return
		} else if args[0] == `search` {
			// Query word
			if len(args) > 1 {
				key := args[1]
				searchAndPrint(key)
				return
			}
		}
	}
	// Default helper
	fmt.Println(`Usage: main [action] [options]`)
}
