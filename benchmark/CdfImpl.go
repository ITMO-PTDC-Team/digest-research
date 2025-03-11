package main

import (
	"bufio"
	"os"
	"sort"
	"strconv"
)

type CDF interface {
	Quantile(v float64) float64
}

type CdfImpl struct {
	data []float64
}

func NewCdfImpl() *CdfImpl {
	return &CdfImpl{
		data: []float64{},
	}
}

func LoadCdfImpl(filename string) (*CdfImpl, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data []float64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		value, err := strconv.ParseFloat(line, 64)
		if err != nil {
			return nil, err
		}
		data = append(data, value)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	sort.Float64s(data)

	return &CdfImpl{data: data}, nil
}

func (cdf *CdfImpl) Quantile(q float64) float64 {
	if len(cdf.data) == 0 {
		return 0
	}
	index := int(q * float64(len(cdf.data)))
	if index >= len(cdf.data) {
		index = len(cdf.data) - 1
	}
	return cdf.data[index]
}

func (cdf *CdfImpl) Cdf(value float64) float64 {
	index := sort.Search(len(cdf.data), func(i int) bool {
		return cdf.data[i] > value
	})
	return float64(index) / float64(len(cdf.data))
}

func CompareCDF(cdf, td CDF, quantiles []float64) []float64 {
	var results []float64
	for _, q := range quantiles {
		cdfValue := cdf.Quantile(q)
		tdValue := td.Quantile(q)
		results = append(results, (tdValue-cdfValue)/cdfValue)
	}
	return results
}

func GetQuantiles(td CDF, quantiles []float64) []float64 {
	var results []float64
	for _, q := range quantiles {
		tdValue := td.Quantile(q)	
		results = append(results, tdValue)
	}
	return results
}

