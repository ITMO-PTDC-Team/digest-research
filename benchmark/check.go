package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"fmt"
)

func (t *TDigest) ProcessedToString()string{
	return t.processed.ListToString()
}
func main() {

	outputFile := "results/CentroidLists.txt"
	outFile, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()
	writer := bufio.NewWriter(outFile)
	
	for i := 0; i < 10; i++ {
		inputFile := "distributions/test_distribution_" + strconv.Itoa(i) + ".txt"
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
		
		
		if(i==0){
		fmt.Fprintf(outFile,"%.4f\t",td.Count())
		fmt.Fprintf(outFile,td.processed.ListToString())
		fmt.Fprintf(outFile,"\n+++++++++++++++++++UNPROCESSED+++++++++++++++++\n")
		fmt.Fprintf(outFile,"%d\n",td.processed.Len())
		}
	}
	fmt.Fprintf(outFile,"Zipf distributions\n")
	writer.Flush()
}