package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"fmt"

)

var quantiles = []float64{0.25, 0.5, 0.75, 0.90, 0.95, 0.99, 0.999}
const NUMBER_OF_DISTRIBUTIONS = 8

func main() {
	td := New();

	inputFile := "input.txt"
	cdf, err := LoadCdfImpl(inputFile)
	if err != nil {
		panic(err)
	}
	
	
	file, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	outputFile := "result.txt"
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
		
	result := CompareCDF(cdf, td, quantiles)

		
	
	for i:=0;i<len(quantiles);i++{
		fmt.Fprintf(outFile,"%.4f\t", quantiles[i])
	}
	fmt.Fprintf(outFile,"\n")
	for i:=0;i<len(quantiles);i++{
		fmt.Fprintf(outFile,"%.4f\t", result[i])
	}
	fmt.Fprintf(outFile,"\n")
		
	writer.Flush()
}