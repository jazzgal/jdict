package main

import (
	"fmt"
	"jdict/lib/jmdict"
)

func main() {
	fmt.Println("Populate Database ...")
	jmdict.PopulateData("data/sample.dict")
}
