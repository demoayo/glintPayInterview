package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	GBP       = "GBP"
	GGM       = "GGM"
	layoutCSV = "02/01/2006 15:04"
)

// ComputeTopNForFilteredSpendersParam ComputeTopNForFilteredSpenders function paramters
type ComputeTopNForFilteredSpendersParam struct {
	// Filters indicates which spenders should be included in computation
	Filters []Filters
	// Spenders specfies list of spenders
	Spenders []*Spender
	// TopNCount specifies the topN count
	TopNCount int
}

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

// Spec specifies computation tasks parameter
type Spec struct {
	FileName  string    `json:"file_name"`
	Filters   []Filters `json:"filters"`
	TopNCount int       `json:"topN"`
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

func main() {
	args := os.Args[1]
	// Tansform program arg to program spec
	spec, err := TransformArgsToSpec(args)
	if err != nil {
		log.Fatal(err)
	}

	// Create new spenders for specified Filename
	spenders, err := CreateNewSpendersFromFile(spec.FileName)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Compute topN for filtered subset of spender
	topNSpenders := ComputeTopNForFilteredSpenders(&ComputeTopNForFilteredSpendersParam{
		Spenders:  spenders,
		Filters:   spec.Filters,
		TopNCount: spec.TopNCount,
	})

	//Print results
	if len(topNSpenders) > 0 {
		// Prints out topN spender to terminal
		PrintOutTopNSpenders(topNSpenders)

		// Format result as JSON
		// Can specify output format, output file name in args, allowing for future format change
		_, err = FormatSpenderToJSONOutput(topNSpenders)
		if err != nil {
			log.Fatal(err)
		}
	}

}

//ComputeTopNForFilteredSpenders returns topN spenders and applies filters
func ComputeTopNForFilteredSpenders(param *ComputeTopNForFilteredSpendersParam) []*Spender {
	topNSpenders := []*Spender{}
	for _, currentSpender := range param.Spenders {
		//Calculate current spenders total spend
		currentSpender.TotalSpend = CalculateTotalSpend(currentSpender)

		qry := &MatchFilterCriteriaParam{
			Filters: param.Filters,
			Spender: currentSpender,
		}
		//Check if current spender match filter criterias
		if !MatchFilterCriteria(qry) {
			continue
		}

		automaticallyFillInitialTopNSpenders := len(topNSpenders) <= (param.TopNCount - 1)
		switch automaticallyFillInitialTopNSpenders {
		case true:
			// Case fills intial topNSpenders
			topNSpenders = append(topNSpenders, currentSpender)
			if len(topNSpenders) >= 1 {
				// Reorder topNSpenders
				sort.Slice(topNSpenders, func(i, j int) bool {
					first, _ := topNSpenders[i].TotalSpend.Float64()
					next, _ := topNSpenders[j].TotalSpend.Float64()
					return first > next
				})
			}
			continue
		case false:
			// Case:: item after the initial topNSpender Slice has been filled
			//TotalSpend calculation Nil Checks
			if topNSpenders[param.TopNCount-1].TotalSpend == nil || currentSpender.TotalSpend == nil {
				continue
			}
			// Check if current spend is greater than last item in the topN spender list
			// If newSpend greater than last item in the topN spender slice, replace item in topN spenders slice
			lastSpenderInTopN, _ := topNSpenders[param.TopNCount-1].TotalSpend.Float64()
			currentTotalSpend, _ := currentSpender.TotalSpend.Float64()
			if currentTotalSpend > lastSpenderInTopN {
				// Swap in current spender to topNSpenders
				topNSpenders[param.TopNCount-1] = currentSpender
				// Reorder topNSpenders
				sort.Slice(topNSpenders, func(i, j int) bool {
					first, _ := topNSpenders[i].TotalSpend.Float64()
					next, _ := topNSpenders[j].TotalSpend.Float64()
					return first > next
				})
			}
		}

	}

	return topNSpenders
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

// CreateNewSpendersFromFile reads a file and returns a list of spenders
func CreateNewSpendersFromFile(fileName string) ([]*Spender, error) {
	// open file
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	//close file
	defer f.Close()

	// read csv values
	csvReader := csv.NewReader(f)
	//Counter used to skip first line
	currentSpenderIndex := 0

	spenders := []*Spender{}
	for {
		// Read a spender at a time
		rec, err := csvReader.Read()
		// Check error
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		// Ignore first line containing table info
		if currentSpenderIndex == 0 {
			currentSpenderIndex++
			continue
		}
		// Create new spender struct
		newSpender, err := CreateNewSpender(rec)
		//Skip spender with error.
		//TODO: Should return error
		if err != nil {
			continue
		}
		//Add current spender to list of spenders
		spenders = append(spenders, newSpender)

		currentSpenderIndex++
	}

	return spenders, nil
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

//TransformArgsToSpecToFilters transforms filter arguments from JSON string to a golang struct
func TransformArgsToSpec(param string) (*Spec, error) {
	var spec *Spec

	err := json.Unmarshal([]byte(param), &spec)
	if err != nil {
		return nil, err
	}

	return spec, nil
}
