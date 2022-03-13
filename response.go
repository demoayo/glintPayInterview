package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

type ComputeTopNForFilteredSpendersResponse struct {
	Spenders []*Spender
}

type TopNSpendersResponse struct {
	Spenders       []*Spender
	OutputFilePath string
}

func FormatSpenderToJSONOutput(spenders []*Spender) (string, error) {
	jsonByte, err := json.MarshalIndent(spenders, " ", " ")

	if err != nil {
		return "", err
	}

	outputFileName := strings.Join([]string{"top_spender", time.Now().Format(time.RFC3339), "json"}, ".")
	ioutil.WriteFile(outputFileName, jsonByte, 0644)

	fmt.Println("\nRelative file path:")
	fmt.Println(outputFileName)

	return outputFileName, err

}

//PrintOutTopNSpenders prints topN spenders
func PrintOutTopNSpenders(spenders []*Spender) {
	fmt.Println("Top spenders:")
	for _, s := range spenders {
		fmt.Printf("Email: %+v, FirstName: %+v, LastName: %+v, MerchantCode: %+v, TotalSpend: %+v \n",
			s.Email,
			s.FirstName,
			s.LastName,
			s.MerchantCode,
			s.TotalSpend,
		)
	}
}

func (resp *TopNSpendersResponse) Apply() {
	//Print response
	if len(resp.Spenders) > 0 {
		// Prints out topN spender to terminal
		PrintOutTopNSpenders(resp.Spenders)

		// Format result as JSON
		// Can specify output format, output file name in args, allowing for future format change
		var err error
		resp.OutputFilePath, err = FormatSpenderToJSONOutput(resp.Spenders)
		if err != nil {
			log.Fatal(err)
		}
	}
}
