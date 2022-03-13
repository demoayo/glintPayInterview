package main

import (
	"math/big"
	"strconv"
	"strings"
	"time"
)

// Filters determines which spenders are to be consider in the computation
type Filters struct {
	Cmp   string `json:"cmp"`
	Field string `json:"field"`
	Value string `json:"value"`
}

// MatchFilterCriteriaParam MatchFilterCriteriaParam function parameters
type MatchFilterCriteriaParam struct {
	Filters []Filters
	Spender *Spender
}

// Spender a BUY/SELL activity by a user
type Spender struct {
	// FirstName spenders first name
	FirstName string
	// LaseName spenders last name
	LastName string
	// Email spenders email address
	Email string
	// Description describes the spend can be CARD SPEND, SELL GOLD, BUY GOLG
	Description string
	// Amount spend quantity
	// TODO: Clarify
	Amount *big.Float
	// FromCurrency spend initial currency can be GBP, GGM
	// GGM stands for cost of gold
	FromCurrency string
	// ToCurrency spend desired currency can be  GBP, GGM
	ToCurrency string
	//Rate an exchange rate??
	// TODO: Clarify
	Rate *big.Float
	//MerchantCode
	MerchantCode string
	Date         time.Time
	Month        int
	//TotalSpend stores the calculation of total spend
	// TODO: Clarify formula assume calculation relates to Amount and Rate
	TotalSpend *big.Float
}

//CreateNewSpender creates and returns a new spender structure
func CreateNewSpender(rec []string) (*Spender, error) {
	var err error
	var amount float64
	if amount, err = strconv.ParseFloat(rec[5], 64); err != nil {
		return nil, err
	}

	var rate float64
	if rate, err = strconv.ParseFloat(rec[8], 64); err != nil {
		return nil, err
	}

	var date time.Time
	if date, err = time.Parse(layoutCSV, rec[9]); err != nil {
		return nil, err
	}

	return &Spender{
		FirstName:    rec[0],
		LastName:     rec[1],
		Email:        rec[2],
		Description:  rec[3],
		MerchantCode: rec[4],
		Amount:       big.NewFloat(amount),
		FromCurrency: rec[6],
		ToCurrency:   rec[7],
		Rate:         big.NewFloat(rate),
		Date:         date,
		Month:        int(date.Month()),
	}, nil
}

//MatchFilterCriteria flag checks if current spenders matches the list of filter
func MatchFilterCriteria(param *MatchFilterCriteriaParam) bool {
	//Filtering is not specified
	if param.Filters == nil || len(param.Filters) == 0 {
		return true
	}

	//Flag checks if spender matches filter
	match := false

	for _, f := range param.Filters {
		processedField := strings.TrimSpace(strings.ToLower(f.Field))
		processedCmp := strings.TrimSpace(strings.ToLower(f.Cmp))
		processedFilterValue := strings.TrimSpace(strings.ToLower(f.Value))

		switch processedField {
		case "description":
			processedSpenderValue := strings.TrimSpace(strings.ToLower(param.Spender.Description))
			//Cmp (description = x)
			if processedCmp == "=" || processedCmp == "==" {
				match = (processedSpenderValue == processedFilterValue)
			}

		case "month":
			processedSpenderValue := param.Spender.Month
			monthFilter, _ := strconv.Atoi(processedFilterValue)
			//Cmp:: (month = x)
			if processedCmp == "=" || processedCmp == "==" {
				match = (processedSpenderValue == monthFilter)
			}
			//Cmp:: (month > x)
			if processedCmp == ">" {
				match = (processedSpenderValue > monthFilter)
			}
			//Cmp:: month >= x
			if processedCmp == ">=" {
				match = (processedSpenderValue >= monthFilter)
			}
			//Cmp:: month < x
			if processedCmp == "<" {
				match = (processedSpenderValue < monthFilter)
			}
			//Cmp:: month <= x
			if processedCmp == "<=" {
				match = (processedSpenderValue <= monthFilter)
			}
		}

		// Return false on encountering any false false
		if !match {
			return false
		}

	}

	return match

}

// CalculateTotalSpend calculates a spenders total spend
// TODO: Verify computation of totalSpend
func CalculateTotalSpend(spender *Spender) *big.Float {
	t := new(big.Float).SetPrec(spender.Amount.Prec())

	if spender.FromCurrency == GBP && spender.ToCurrency == GBP {
		return spender.Amount
	} else if spender.FromCurrency == GBP && spender.ToCurrency == GGM {
		return spender.Amount
	} else {
		return t.Mul(spender.Amount, spender.Rate)
	}

}
