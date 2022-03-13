package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"sort"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

//ComputeTopNForFilteredSpenders returns topN spenders and applies filters
func (s *Service) ComputeTopNForFilteredSpenders(req *ComputeTopNForFilteredSpendersRequest) *ComputeTopNForFilteredSpendersResponse {
	topNSpenders := []*Spender{}
	for _, currentSpender := range req.Spenders {
		//Calculate current spenders total spend
		currentSpender.TotalSpend = CalculateTotalSpend(currentSpender)

		qry := &MatchFilterCriteriaParam{
			Filters: req.Filters,
			Spender: currentSpender,
		}
		//Check if current spender match filter criterias
		if !MatchFilterCriteria(qry) {
			continue
		}

		automaticallyFillInitialTopNSpenders := len(topNSpenders) <= (req.TopNCount - 1)
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
			if topNSpenders[req.TopNCount-1].TotalSpend == nil || currentSpender.TotalSpend == nil {
				continue
			}
			// Check if current spend is greater than last item in the topN spender list
			// If newSpend greater than last item in the topN spender slice, replace item in topN spenders slice
			lastSpenderInTopN, _ := topNSpenders[req.TopNCount-1].TotalSpend.Float64()
			currentTotalSpend, _ := currentSpender.TotalSpend.Float64()
			if currentTotalSpend > lastSpenderInTopN {
				// Swap in current spender to topNSpenders
				topNSpenders[req.TopNCount-1] = currentSpender
				// Reorder topNSpenders
				sort.Slice(topNSpenders, func(i, j int) bool {
					first, _ := topNSpenders[i].TotalSpend.Float64()
					next, _ := topNSpenders[j].TotalSpend.Float64()
					return first > next
				})
			}
		}

	}

	return &ComputeTopNForFilteredSpendersResponse{
		Spenders: topNSpenders,
	}
}

// CreateNewSpendersFromFile reads a file and returns a list of spenders
func (s *Service) CreateNewSpendersFromFile(req *CreateNewSpendersFromFileRequest) ([]*Spender, error) {
	// open file
	f, err := os.Open(req.FileName)
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

func (s *Service) TopNSpenders(req *TopNSpendersRequest) *TopNSpendersResponse {
	// Create new spenders for specified Filename
	spenders, err := s.CreateNewSpendersFromFile(&CreateNewSpendersFromFileRequest{
		FileName: req.FileName,
	})
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Compute topN for filtered subset of spender
	computeTopNForFilteredSpendersRes := s.ComputeTopNForFilteredSpenders(&ComputeTopNForFilteredSpendersRequest{
		Spenders:  spenders,
		Filters:   req.Filters,
		TopNCount: req.TopNCount,
	})

	// Apply response
	topNSpendersResponse := &TopNSpendersResponse{
		Spenders: computeTopNForFilteredSpendersRes.Spenders,
	}
	topNSpendersResponse.Apply()

	return topNSpendersResponse

}
