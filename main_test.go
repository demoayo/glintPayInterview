package main

import (
	"log"
	"math/big"
	"testing"
)

func TestArgsProcessing(t *testing.T) {
	testCases := []struct {
		name   string
		desc   string
		expect bool
		params string
		actual func(params string) bool
	}{
		{
			name: "main.TransformArgs",
			desc: "Generate spec",
			params: `{
				"file_name": "sample-transactions.csv", 
			    "filters": [
					{"field": "description", "cmp": "=", "value": "CARD SPEND"},
					{"field": "month", "cmp": "=", "value": "1"}
					], 
				"top_n": 5}`,
			expect: true,
			actual: func(params string) bool {
				req, err := TransformArgsToTopNSpendersRequest(params)
				if err != nil {
					log.Println(err.Error())
				}

				return req != nil && len(req.Filters) == 2 && req.TopNCount == 5

			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actual := tC.actual(tC.params)
			if actual != tC.expect {
				t.Errorf("\nExpecting: %v \nActual: %v \nDescription: %v \n", tC.expect, actual, tC.desc)
			}
		})
	}
}

func TestMatchFilterCriteria(t *testing.T) {
	testCases := []struct {
		name   string
		desc   string
		expect bool
		params *MatchFilterCriteriaParam
		actual func(params *MatchFilterCriteriaParam) bool
	}{
		{
			name: "main.TestMatchFilterCriteria",
			desc: "Filter matches CARD SPEND",
			params: &MatchFilterCriteriaParam{
				Filters: []Filters{
					{Field: "Description", Cmp: "=", Value: "CARD SPEND"},
				},
				Spender: &Spender{Description: "CARD SPEND"},
			},
			expect: true,
			actual: func(params *MatchFilterCriteriaParam) bool {
				return MatchFilterCriteria(params)
			},
		},
		{
			name: "main.TestMatchFilterCriteria",
			desc: "Filter matches CARD SPEND AND MONTH",
			params: &MatchFilterCriteriaParam{
				Filters: []Filters{
					{Field: "Description", Cmp: "=", Value: "CARD SPEND"},
					{Field: "Month", Cmp: "=>", Value: "1"},
				},
				Spender: &Spender{Description: "CARD SPEND", Month: 2},
			},
			expect: true,
			actual: func(params *MatchFilterCriteriaParam) bool {
				return MatchFilterCriteria(params)
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actual := tC.actual(tC.params)
			if actual != tC.expect {
				t.Errorf("\nExpecting: %v \nActual: %v \nDescription: %v \n", tC.expect, actual, tC.desc)
			}
		})
	}
}

func TestComputeTopNForFilteredSpenders(t *testing.T) {
	testCases := []struct {
		name   string
		desc   string
		expect bool
		params *ComputeTopNForFilteredSpendersRequest
		actual func(params *ComputeTopNForFilteredSpendersRequest) bool
	}{
		{
			name: "service.ComputeTopNForFilteredSpenders",
			desc: "Calculates Top N for filtered Spenders using service",
			params: &ComputeTopNForFilteredSpendersRequest{
				TopNCount: 10,
				Filters: []Filters{
					{Field: "Description", Cmp: "=", Value: "CARD SPEND"},
					{Field: "Month", Cmp: "=", Value: "2"},
				},
				Spenders: []*Spender{
					{Email: "2@email.com", Description: "CARD SPEND", Amount: big.NewFloat(100.00), FromCurrency: GBP, ToCurrency: GBP, Rate: big.NewFloat(1.00), Month: 2},
					{Email: "1@email.com", Description: "CARD SPEND", Amount: big.NewFloat(200.00), FromCurrency: GBP, ToCurrency: GBP, Rate: big.NewFloat(1.00), Month: 2},
					{Email: "4@email.com", Description: "SELL GOLD", Amount: big.NewFloat(200.00), FromCurrency: GBP, ToCurrency: GBP, Rate: big.NewFloat(1.00), Month: 2},
					{Email: "3@email.com", Description: "CARD SPEND", Amount: big.NewFloat(50.00), FromCurrency: GBP, ToCurrency: GGM, Rate: big.NewFloat(1.00), Month: 2},
					{Email: "wrongMonth@email.com", Description: "CARD SPEND", Amount: big.NewFloat(50.00), FromCurrency: GBP, ToCurrency: GGM, Rate: big.NewFloat(1.00), Month: 3},
				},
			},
			expect: true,
			actual: func(params *ComputeTopNForFilteredSpendersRequest) bool {
				service := NewService()
				resp := service.ComputeTopNForFilteredSpenders(params)
				return (resp.Spenders != nil &&
					(len(resp.Spenders) == 3) &&
					(resp.Spenders[0].Email == "1@email.com") &&
					(resp.Spenders[1].Email == "2@email.com") &&
					(resp.Spenders[2].Email == "3@email.com"))
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actual := tC.actual(tC.params)
			if actual != tC.expect {
				t.Errorf("\nExpecting: %v \nActual: %v \nDescription: %v \n", tC.expect, actual, tC.desc)
			}
		})
	}
}

func TestTopNSpenders(t *testing.T) {
	testCases := []struct {
		name   string
		desc   string
		expect bool
		params *TopNSpendersRequest
		actual func(params *TopNSpendersRequest) bool
	}{
		{
			name: "service.ComputeTopNForFilteredSpenders",
			desc: "Calculates Top 5 Spenders using service",
			params: &TopNSpendersRequest{
				TopNCount: 5,
				Filters: []Filters{
					{Field: "Description", Cmp: "=", Value: "CARD SPEND"},
					{Field: "Month", Cmp: "=", Value: "2"},
				},
				FileName: "sample-transactions.csv",
			},
			expect: true,
			actual: func(params *TopNSpendersRequest) bool {
				service := NewService()
				resp := service.TopNSpenders(params)
				return (resp.Spenders != nil &&
					(resp.OutputFilePath != "") &&
					(len(resp.Spenders) == 5) &&
					(resp.Spenders[0].Email == "kaif.beck@mailinator.com") &&
					(resp.Spenders[1].Email == "ceara.valdez@mailinator.com") &&
					(resp.Spenders[2].Email == "cosmo.mansell@mailinator.com") &&
					(resp.Spenders[3].Email == "andy.nguyen@mailinator.com") &&
					(resp.Spenders[4].Email == "finlay.rasmussen@mailinator.com"))

			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actual := tC.actual(tC.params)
			if actual != tC.expect {
				t.Errorf("\nExpecting: %v \nActual: %v \nDescription: %v \n", tC.expect, actual, tC.desc)
			}
		})
	}
}
