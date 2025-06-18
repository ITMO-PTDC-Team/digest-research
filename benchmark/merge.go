package main

import (
	"fmt"
	"log"
	"os"
	"bufio"
	"strconv"
	"strings"
	"sort"
)
var quantiles = []float64{}

func main() {
	cdfMinutes := make([]*CdfImpl, 0)
	tdMinutes := make([]*TDigest, 0)

	for i:=0.5 ; i <1 ; i+=0.001{
	
		quantiles = append(quantiles, i)
	}
	sort.Float64s(quantiles)

	results :=make ([][]float64,60)
	for i := range results{
		results[i]=make([]float64, len(quantiles))
	}

	outputFile := "nginx_request_length_res_td_raw.txt"
	outFile, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()
	writer := bufio.NewWriter(outFile)

	for start := 1; start <= 3600; start += 60 {
		end := start + 59
		if end > 3600 {
			end = 3600
		}

		minuteCdf := NewCdfImpl()
		minuteTd := NewWithCompression(100)
		for i := start; i <= end; i++ {
			filename := fmt.Sprintf("generated/nginx_request_length-1s_%d.txt", i)
			cdf, err := LoadCdfImpl(filename)
			if err != nil {
				log.Printf("Ошибка загрузки файла %s: %v", filename, err)
				continue
			}

			file, err := os.Open(filename)
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
				minuteTd.Add(value, 1)
			}

			if err := scanner.Err(); err != nil {
				panic(err)
			}
			
			result := CompareCDF(cdf, minuteTd, quantiles)
			for j := range result {
				results[start/60][j]=result[j]
			}
			

			minuteCdf.Merge(cdf)
		}
		tdMinutes = append(tdMinutes, minuteTd)
		cdfMinutes = append(cdfMinutes, minuteCdf)
	}

	resultCdf := NewCdfImpl()
	resultCdf.MergeArray(cdfMinutes)
	resultTd := NewWithCompression(100)
	for _, td := range tdMinutes {
        resultTd.Merge(td)
    }

	
	// result := CompareCDF(resultCdf, resultTd, quantiles)


	// fmt.Fprintf(outFile,"%d %d\n",len(quantiles),60+2)
	// for i:=0;i<len(quantiles);i++{
	// 	fmt.Fprintf(outFile,"%.4f\t", quantiles[i])
	// }
	// fmt.Fprintf(outFile,"\n")

	// for i:=0;i<len(quantiles);i++{
	// 	fmt.Fprintf(outFile,"%.4f\t", result[i])
	// }
	// fmt.Fprintf(outFile,"\n")

	for i:=0;i<60;i++{
		for j:=0;j<len(quantiles);j++{
			fmt.Fprintf(outFile,"%.4f\n", results[i][j])
		}
		// fmt.Fprintf(outFile,"\n")
	}
	writer.Flush()
}
