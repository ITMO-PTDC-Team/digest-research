package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"fmt"
	"math/rand"
	"sort"
)

var quantiles = []float64{0.25, 0.5, 0.75,0.80, 0.90}
const NUMBER_OF_DISTRIBUTIONS = 10
const NUMBER_OF_ZIPF_DISTRIBUTIONS=5
const NUMBER_OF_NORMAL_DISTRIBUTIONS=5

func main() {
	rand.Seed(42)
	for i:=0 ; i < 5 ; i++{
		randomNum := 0.9 + rand.Float64()*(1.0-0.9)
		quantiles = append(quantiles, randomNum)
	}
	for i:=0 ; i < 5 ; i++{
		randomNum := 0.99 + rand.Float64()*(1.0-0.99)
		quantiles = append(quantiles, randomNum)
	}
	sort.Float64s(quantiles)
	results :=make ([][]float64,NUMBER_OF_DISTRIBUTIONS)
	for i := range results{
		results[i]=make([]float64, len(quantiles))
	}
	outputFile := "results/resultPOW.txt"
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
	fmt.Fprintf(outFile,"Zipf distributions\n")
	fmt.Fprintf(outFile,"\n")
	for i:=0;i<len(quantiles);i++{
		fmt.Fprintf(outFile,"%.4f\t", quantiles[i])
	}
	fmt.Fprintf(outFile,"\n")
	fmt.Fprintf(outFile,"\n")
		for i := 0; i < NUMBER_OF_ZIPF_DISTRIBUTIONS; i++ {
			for j := 0 ; j < len(quantiles); j++{
			fmt.Fprintf(outFile,"%.4f\t", results[i][j])
		}
		fmt.Fprintf(outFile,"\n")
	}
	fmt.Fprintf(outFile,"Normal distributions\n")
	fmt.Fprintf(outFile,"\n")
	for i:=0;i<len(quantiles);i++{
		fmt.Fprintf(outFile,"%.4f\t", quantiles[i])
	}
	fmt.Fprintf(outFile,"\n")
	fmt.Fprintf(outFile,"\n")
		for i := 0; i < NUMBER_OF_NORMAL_DISTRIBUTIONS; i++ {
			for j := 0 ; j < len(quantiles); j++{
			fmt.Fprintf(outFile,"%.4f\t", results[i+NUMBER_OF_ZIPF_DISTRIBUTIONS][j])
		}
		fmt.Fprintf(outFile,"\n")
	}
	avg_results:=make([]float64, len(quantiles))
	fmt.Fprintf(outFile,"\nAverage diff\n\n")
	for j := 0 ; j < len(quantiles); j++{
		for i := 0; i < NUMBER_OF_DISTRIBUTIONS; i++ {
			avg_results[j]+=results[i][j]
		}
		avg_results[j]=avg_results[j]/float64(NUMBER_OF_DISTRIBUTIONS)
		fmt.Fprintf(outFile,"%.4f\t", avg_results[j])
	}
	fmt.Fprintf(outFile,"\n")
	writer.Flush()
}