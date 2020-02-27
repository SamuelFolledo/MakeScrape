package main

/*
NOTE:
- $ go test //to run tests
- $ go test -cover //to check test coverage
- $ go tool cover -html=coverage.out //generate coverage.out file which is used to generate a HTML page which shows exactly what lines have been covered
*/

import (
	"testing"
)

func Test_createDealFromScrapedText(t *testing.T) {
	var deal = createDealFromScrapedText("")
	if deal.Name != "" { //check if name is still empty
		t.Error("Expected empty Deal")
	}
	if deal.CurrentPrice != 0 { //check if name is still empty
		t.Error("Expected empty Deal")
	}
	if deal.PreviousPrice != 0 { //check if name is still empty
		t.Error("Expected empty Deal")
	}
}

func TestCalculate(t *testing.T) {
	if Calculate(2) != 4 {
		t.Error("Expected 2 + 2 to equal 4")
	}
}

func TestTableCalculate(t *testing.T) {
	var tests = []struct {
		input    int
		expected int
	}{
		{2, 4},
		{-1, 1},
		{0, 2},
		{-5, -3},
		{99999, 100001},
	}

	for _, test := range tests {
		if output := Calculate(test.input); output != test.expected {
			t.Error("Test Failed: {} inputted, {} expected, recieved: {}", test.input, test.expected, output)
		}
	}
}
