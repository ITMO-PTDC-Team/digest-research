package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"github.com/influxdata/tdigest"
)

var quantiles = []float64{0.25, 0.5, 0.75, 0.90, 0.95, 0.99, 0.999}

func main() {
	cdf, err := LoadCdfImpl("distribution.txt")
	if err != nil {
		panic(err)
	}
	td := tdigest.New()
	file, err := os.Open("distribution.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		value, err := strconv.ParseFloat(strings.TrimSpace(line), 64)
		if err != nil {
			panic(err)
		}
		td.Add(value, 1)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	results := CompareCDF(cdf, td, quantiles)

	outputFile, err := os.Create("result.txt")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	for _, result := range results {
		_, err := writer.WriteString(result + "\n")
		if err != nil {
			panic(err)
		}
	}
	writer.Flush()

}