package main

import (
	"log"
	"os"
)

const (
	GBP       = "GBP"
	GGM       = "GGM"
	layoutCSV = "02/01/2006 15:04"
)

func main() {
	args := os.Args[1]
	// Tansform program arg to program spec
	spec, err := TransformArgsToTopNSpendersRequest(args)
	if err != nil {
		log.Fatal(err)
	}

	// Create new spenders for specified Filename
	service := NewService()
	service.TopNSpenders(spec)

}
