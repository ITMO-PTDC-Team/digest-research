package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"fmt"
	//"github.com/influxdata/tdigest"
	//"github.com/influxdata/tdigest"
)

var quantiles = []float64{0.25, 0.5, 0.75, 0.90, 0.95, 0.99, 0.999, 0.9999}
const NUMBER_OF_DISTRIBUTIONS = 8

func main() {

	results :=make ([][]float64,NUMBER_OF_DISTRIBUTIONS)
	for i := range results{
		results[i]=make([]float64, len(quantiles))
	}
	outputFile := "results/resultEXP.txt"
	outFile, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()
	writer := bufio.NewWriter(outFile)
	
	for i := 0; i < NUMBER_OF_DISTRIBUTIONS; i++ {
		inputFile := "distributions/test_distribution_" + strconv.Itoa(i) + ".txt"
		cdf, err := LoadCdfImpl(inputFile)
		if err != nil {
			panic(err)
		}
		
		td := New();
		
		file, err := os.Open(inputFile)
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
		
		result := CompareCDF(cdf, td, quantiles)
		for j := range result {
			results[i][j]=result[j]
		}
		
	}
	for i:=0;i<len(quantiles);i++{
		fmt.Fprintf(outFile,"%.4f\t", quantiles[i])
	}
	fmt.Fprintf(outFile,"\n")
		for i := 0; i < NUMBER_OF_DISTRIBUTIONS; i++ {
			for j := 0 ; j < len(quantiles); j++{
			fmt.Fprintf(outFile,"%.4f\t", results[i][j])
		}
		fmt.Fprintf(outFile,"\n")
	}
	writer.Flush()
}