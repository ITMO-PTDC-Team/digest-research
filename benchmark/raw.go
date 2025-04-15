package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"fmt"
	"sort"
)

var quantiles = []float64{}
const NUMBER_OF_ZIPF_DISTRIBUTIONS=5
const NUMBER_OF_NORMAL_DISTRIBUTIONS=8
const NUMBER_OF_PARETO_DISTRIBUTIONS=5
const NUMBER_OF_N_NORMAL_DISTIBUTIONS=8
const NUMBER_OF_NORMAL_TAIL_DISTRIBUTION=8
const NUMBER_OF_DISTRIBUTIONS = NUMBER_OF_NORMAL_DISTRIBUTIONS+NUMBER_OF_PARETO_DISTRIBUTIONS+NUMBER_OF_ZIPF_DISTRIBUTIONS+NUMBER_OF_N_NORMAL_DISTIBUTIONS+NUMBER_OF_NORMAL_TAIL_DISTRIBUTION

func main() {

	for i:=0.001 ; i <1 ; i+=0.001{
	
		quantiles = append(quantiles, i)
	}
	sort.Float64s(quantiles)
	results_td :=make ([][]float64,NUMBER_OF_DISTRIBUTIONS)
	for i := range results_td{
		results_td[i]=make([]float64, len(quantiles))
	}
	results_cdf :=make ([][]float64,NUMBER_OF_DISTRIBUTIONS)
	for i := range results_cdf{
		results_cdf[i]=make([]float64, len(quantiles))
	}
	file_name := "sin_pow_2_5"
	outputFile_td := "quantiles/td_" + file_name + ".txt"

	outFile_td, err := os.Create(outputFile_td)
	if err != nil {
		panic(err)
	}
	defer outFile_td.Close()
	writer_td := bufio.NewWriter(outFile_td)
	

	outputFile_cdf := "quantiles/cdf_" + file_name + ".txt"

	outFile_cdf, err := os.Create(outputFile_cdf)
	if err != nil {
		panic(err)
	}
	defer outFile_td.Close()
	writer_cdf := bufio.NewWriter(outFile_cdf)

	for i := 0; i < NUMBER_OF_DISTRIBUTIONS; i++ {
		td := NewWithCompression(100);
		inputFile := "distributions/test_distribution_" + strconv.Itoa(i) + ".txt"
		cdf, err := LoadCdfImpl(inputFile)
		if err != nil {
			panic(err)
		}
		
		
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
		
		result_td := GetQuantiles(td, quantiles)
		for j := range result_td {
			results_td[i][j]=result_td[j]
		}
		result_cdf := GetQuantiles(cdf, quantiles)
		for j := range result_cdf {
			results_cdf[i][j]=result_cdf[j]
		}
		
	}

	fmt.Fprintf(outFile_td,"%d %d\n",len(quantiles),NUMBER_OF_DISTRIBUTIONS+1)
	fmt.Fprintf(outFile_cdf,"%d %d\n",len(quantiles),NUMBER_OF_DISTRIBUTIONS+1)
	for i:=0;i<len(quantiles);i++{
		fmt.Fprintf(outFile_td,"%.4f\t", quantiles[i])
		fmt.Fprintf(outFile_cdf,"%.4f\t", quantiles[i])
	}
	fmt.Fprintf(outFile_td,"\n")
	fmt.Fprintf(outFile_cdf,"\n")
	for i:=0;i<NUMBER_OF_DISTRIBUTIONS;i++{
		for j:=0;j<len(quantiles);j++{
			fmt.Fprintf(outFile_td,"%.4f\t", results_td[i][j])
		}
		fmt.Fprintf(outFile_td,"\n")
		for j:=0;j<len(quantiles);j++{
			fmt.Fprintf(outFile_cdf,"%.4f\t", results_cdf[i][j])
		}
		fmt.Fprintf(outFile_cdf,"\n")
	}

	fmt.Fprintf(outFile_td,"\n")
	
	fmt.Fprintf(outFile_cdf,"\n")
	writer_td.Flush()
	writer_cdf.Flush()

}