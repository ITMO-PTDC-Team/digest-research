package main

import (
	"bufio"
	"os"
	"strconv"
	"fmt"
	"sort"
	"math"
)

var quantiles = []float64{}
const NUMBER_OF_ZIPF_DISTRIBUTIONS=5
const NUMBER_OF_NORMAL_DISTRIBUTIONS=8
const NUMBER_OF_PARETO_DISTRIBUTIONS=5
const NUMBER_OF_N_NORMAL_DISTIBUTIONS=5
const NUMBER_OF_DISTRIBUTIONS = NUMBER_OF_NORMAL_DISTRIBUTIONS+NUMBER_OF_PARETO_DISTRIBUTIONS+NUMBER_OF_ZIPF_DISTRIBUTIONS+NUMBER_OF_N_NORMAL_DISTIBUTIONS

func main() {

	for i:=0.001 ; i <1 ; i+=0.001{
	
		quantiles = append(quantiles, i)
	}
	sort.Float64s(quantiles)
	results :=make ([][]float64,NUMBER_OF_DISTRIBUTIONS)
	for i := range results{
		results[i]=make([]float64, len(quantiles))
	}

	outputFile := "quantiles/cdf100sin_pow_2_5.txt"
	outFile, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()
	writer := bufio.NewWriter(outFile)
	
	for i := 0; i < NUMBER_OF_DISTRIBUTIONS; i++ {
		//td := NewWithCompression(100);
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
		
		// for scanner.Scan() {
		// 	line := scanner.Text()
		// 	value, err := strconv.ParseFloat(strings.TrimSpace(line), 64)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	td.Add(value, 1)
		// }

		if err := scanner.Err(); err != nil {
			panic(err)
		}
		
		result := GetQuantiles(cdf, quantiles)
		for j := range result {
			results[i][j]=result[j]
		}
		
	}

	max_diff:=make([]float64, len(quantiles))
	avg_diff:=make([]float64, len(quantiles))
	//max_diff_q:=0.0
	// max_diff_v:=0.0
	// avg_max_diff_v:=0.0
	//avg_max_diff_q:=0.0
	for j:=0;j<NUMBER_OF_DISTRIBUTIONS;j++{
		for i:=0;i<len(quantiles);i++{
			
			if j==0 {
				max_diff[i]=math.Abs(results[0][i])
				avg_diff[i]=math.Abs(results[0][i])
				}else{
					max_diff[i]=max(math.Abs(results[j][i]),max_diff[i])
				avg_diff[i]+=math.Abs(results[j][i])
			}
			if j==NUMBER_OF_DISTRIBUTIONS-1{
				avg_diff[i]=avg_diff[i]/NUMBER_OF_DISTRIBUTIONS
			}
		}
	}
	
	//fmt.Fprintf(outFile,"MaxDiff\n")
	//fmt.Fprintf(outFile,"\n")
	fmt.Fprintf(outFile,"%d %d\n",len(quantiles),NUMBER_OF_DISTRIBUTIONS+3)
	for i:=0;i<len(quantiles);i++{
		fmt.Fprintf(outFile,"%.4f\t", quantiles[i])
	}
	fmt.Fprintf(outFile,"\n")
	//fmt.Fprintf(outFile,"\n")
	for i:=0;i<NUMBER_OF_DISTRIBUTIONS;i++{
		for j:=0;j<len(quantiles);j++{
			fmt.Fprintf(outFile,"%.4f\t", results[i][j])
		}
		fmt.Fprintf(outFile,"\n")
	}

	//fmt.Fprintf(outFile,"\n")
	//fmt.Fprintf(outFile,"\n")
	
	for j := 0 ; j < len(quantiles); j++{
		fmt.Fprintf(outFile,"%.4f\t", max_diff[j])
		// if max_diff_v<max_diff[j]{
		// 	max_diff_v=max_diff[j]
		// 	max_diff_q=quantiles[j]
		// }
	}
	
	//fmt.Fprintf(outFile,"\n\nMax Diff Value And Quantile\n")
	//fmt.Fprintf(outFile,"Quantile:%.4f Value:%.4f\n",max_diff_q,max_diff_v)
	
	//fmt.Fprintf(outFile,"\nAvg Diff\n")
	//fmt.Fprintf(outFile,"\n")
	
	// for i:=0;i<len(quantiles);i++{
	// 	fmt.Fprintf(outFile,"%.4f\t", quantiles[i])
	// }
	
	// fmt.Fprintf(outFile,"\n")
	fmt.Fprintf(outFile,"\n")
	
	for j := 0 ; j < len(quantiles); j++{
			fmt.Fprintf(outFile,"%.4f\t", avg_diff[j])
			// if avg_max_diff_v<avg_diff[j]{
			// 	avg_max_diff_v=avg_diff[j]
			// 	avg_max_diff_q=quantiles[j]
			// }
		}
		//fmt.Fprintf(outFile,"\n\nMax AVG Diff Value And Quantile\n")
		//fmt.Fprintf(outFile,"Quantile:%.4f Value:%.4f\n",avg_max_diff_q,avg_max_diff_v)
	fmt.Fprintf(outFile,"\n")
	writer.Flush()

}