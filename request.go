package main

import "encoding/json"

// ComputeTopNForFilteredSpendersRequest specifies request structure
type ComputeTopNForFilteredSpendersRequest struct {
	// Filters indicates which spenders should be included in computation
	Filters []Filters
	// Spenders specfies list of spenders
	Spenders []*Spender
	// TopNCount specifies the topN count
	TopNCount int
}

// TopNSpendersRequest specifies request structure
type TopNSpendersRequest struct {
	FileName  string    `json:"file_name"`
	Filters   []Filters `json:"filters"`
	TopNCount int       `json:"top_n"`
}

//CreateNewSpendersFromFileReq specifies co
type CreateNewSpendersFromFileRequest struct {
	FileName string
}

//TransformArgsToTopNSpendersRequestToFilters transforms filter arguments from JSON string to a golang struct
func TransformArgsToTopNSpendersRequest(args string) (*TopNSpendersRequest, error) {
	var req *TopNSpendersRequest

	err := json.Unmarshal([]byte(args), &req)
	if err != nil {
		return nil, err
	}

	return req, nil
}
