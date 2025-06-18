package main

import (
	"bufio"
	"os"

	"fmt"
	"math"
)

var quantiles = []float64{0.5, 0.75, 0.8, 0.85, 0.9, 0.95, 0.99, 0.995, 0.999}

func main() {

	results := make([][]float64, 60)
	for i := range results {
		results[i] = make([]float64, len(quantiles))
	}
	filename := "lentach_grpc_time_res_td_sqrt_raw"
	outputFile := "res/" + filename + ".txt"
	outFile, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()
	writer := bufio.NewWriter(outFile)

	for i := 0; i < 60; i++ {
		//td := NewWithCompression(100);
		inputFile := filename + ".txt"
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
			results[i][j] = result[j]
		}

	}

	max_diff := make([]float64, len(quantiles))
	avg_diff := make([]float64, len(quantiles))
	max_diff_q := 0.0
	max_diff_v := 0.0
	// avg_max_diff_v:=0.0
	//avg_max_diff_q:=0.0
	for j := 0; j < 60; j++ {
		for i := 0; i < len(quantiles); i++ {

			if j == 0 {
				max_diff[i] = math.Abs(results[0][i])
				avg_diff[i] = math.Abs(results[0][i])
			} else {
				max_diff[i] = max(math.Abs(results[j][i]), max_diff[i])
				avg_diff[i] += math.Abs(results[j][i])
			}
			if j == 60-1 {
				avg_diff[i] = avg_diff[i] / 60
			}
		}
	}

	//fmt.Fprintf(outFile,"MaxDiff\n")
	//fmt.Fprintf(outFile,"\n")
	fmt.Fprintf(outFile, "%d %d\n", len(quantiles), 60+3)
	for i := 0; i < len(quantiles); i++ {
		fmt.Fprintf(outFile, "%.4f\t", quantiles[i])
	}
	fmt.Fprintf(outFile, "\n")
	//fmt.Fprintf(outFile,"\n")
	for i := 0; i < 60; i++ {
		for j := 0; j < len(quantiles); j++ {
			fmt.Fprintf(outFile, "%.4f\t", results[i][j])
		}
		fmt.Fprintf(outFile, "\n")
	}

	//fmt.Fprintf(outFile,"\n")
	//fmt.Fprintf(outFile,"\n")

	for j := 0; j < len(quantiles); j++ {
		fmt.Fprintf(outFile, "%.4f\t", max_diff[j])
		if max_diff_v < max_diff[j] {
			max_diff_v = max_diff[j]
			max_diff_q = quantiles[j]
		}
	}

	fmt.Fprintf(outFile, "\n\nMax Diff Value And Quantile\n")
	fmt.Fprintf(outFile, "Quantile:%.4f Value:%.4f\n", max_diff_q, max_diff_v)

	//fmt.Fprintf(outFile,"\nAvg Diff\n")
	//fmt.Fprintf(outFile,"\n")

	// for i:=0;i<len(quantiles);i++{
	// 	fmt.Fprintf(outFile,"%.4f\t", quantiles[i])
	// }

	// fmt.Fprintf(outFile,"\n")
	//fmt.Fprintf(outFile,"\n")

	for j := 0; j < len(quantiles); j++ {
		//fmt.Fprintf(outFile,"%.4f\t", avg_diff[j])
		// if avg_max_diff_v<avg_diff[j]{
		// 	avg_max_diff_v=avg_diff[j]
		// 	avg_max_diff_q=quantiles[j]
		// }
	}
	//fmt.Fprintf(outFile,"\n\nMax AVG Diff Value And Quantile\n")
	//fmt.Fprintf(outFile,"Quantile:%.4f Value:%.4f\n",avg_max_diff_q,avg_max_diff_v)
	//fmt.Fprintf(outFile,"\n")
	writer.Flush()

}
