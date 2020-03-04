package main

/*
NOTE:
- $ go test //to run tests
UNIT TEST
- $ go test -cover //to check test coverage
- $ go tool cover -html=coverage.out //generate coverage.out file which is used to generate a HTML page which shows exactly what lines have been covered
BENCHMARK TEST
- $ go test -bench=. //runs all benchmarks within our package
- $ go test -run=Calculate -bench=. //runs all bench test with Calculate in test function name

TEST GENERATOR: https://github.com/cweill/gotests
*/

import (
	"testing"
)

func Test_createDealFromScrapedText(t *testing.T) {
	var deal = createDealFromScrapedText("") //simple tests
	if deal.Name != "" {                     //check if name is still empty
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

func TestTableCalculate(t *testing.T) { //T is a type passed to Test functions to manage test state and support formatted test logs
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

func benchmarkCalculate(input int, b *testing.B) { //B is a type passed to Benchmark functions to manage benchmark timing and to specify the number of iterations to run
	for n := 0; n < b.N; n++ {
		Calculate(input)
	}
}

func BenchmarkCalculate100(b *testing.B)         { benchmarkCalculate(100, b) }
func BenchmarkCalculateNegative100(b *testing.B) { benchmarkCalculate(-100, b) }
func BenchmarkCalculateNegative1(b *testing.B)   { benchmarkCalculate(-1, b) }
