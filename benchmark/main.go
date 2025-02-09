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
	for i := 0; i < 5; i++ {
		inputFile := "test_distribution_" + strconv.Itoa(i) + ".txt"
		outputFile := "result_" + strconv.Itoa(i) + ".txt"
		cdf, err := LoadCdfImpl(inputFile)
		if err != nil {
			panic(err)
		}

		td := tdigest.New()

		file, err := os.Open(inputFile)
		if err != nil {
			panic(err)
		}
		defer file.Close()


		outFile, err := os.Create(outputFile)
		if err != nil {
			panic(err)
		}
		defer outFile.Close()

		writer := bufio.NewWriter(outFile)

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
		for _, result := range results {
			_, err := writer.WriteString(result + "\n")
			if err != nil {
				panic(err)
			}
		}

		writer.Flush()
	}
}