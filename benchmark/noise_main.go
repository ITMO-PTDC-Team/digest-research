package main

import (
    "bufio"
    "fmt"
    "os"
    "path/filepath"
    "strconv"
    "strings"
)

var quantiles = []float64{0.25, 0.5, 0.75, 0.80, 0.90, 0.91,0.92,0.93,0.94,0.95,0.96,0.97,0.98,0.99}

func loadTDigestFromFile(filename string) *TDigest {
    file, err := os.Open(filename)
    if err != nil {
        panic(fmt.Sprintf("ошибка открытия файла %s: %v", filename, err))
    }
    defer file.Close()

    td := New()
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" {
            continue
        }
        value, err := strconv.ParseFloat(line, 64)
        if err != nil {
            panic(fmt.Sprintf("ошибка преобразования строки '%s' в число из файла %s: %v", line, filename, err))
        }
        td.Add(value, 1)
    }
    if err := scanner.Err(); err != nil {
        panic(err)
    }
    return td
}

func main() {
    inputFolder := "noise"
    outputFolder := filepath.Join("result_noise", "result")
    os.MkdirAll(outputFolder, 0755)
    outputFile := filepath.Join(outputFolder, "merge_1000_2.txt")
    
    outFile, err := os.Create(outputFile)
    if err != nil {
        panic(err)
    }
    defer outFile.Close()
    writer := bufio.NewWriter(outFile)

    globalTD := New()

    for group := 1; group <= 60; group++ {
        groupTD := New()
        for sub := 1; sub <= 60; sub++ {
            filename := filepath.Join(inputFolder, fmt.Sprintf("merge_%d_%d.txt", group, sub))
            td := loadTDigestFromFile(filename)
            groupTD.Merge(td)
            groupTD.updateCumulative()
        }
    }

    noiseFile := filepath.Join(inputFolder, "merge_noise.txt")
    noiseTD := loadTDigestFromFile(noiseFile)
    globalTD.Merge(noiseTD)
    globalTD.updateCumulative()
    totalFile := filepath.Join(inputFolder, "total.txt")
    realCdf, err := LoadCdfImpl(totalFile)
    if err != nil {
        panic(err)
    }

    diffs := CompareCDF(realCdf, globalTD, quantiles)

    writer.WriteString("\nСравнение квантилей (TDigest vs настоящие):\n")
    for i, q := range quantiles {
        realVal := realCdf.Quantile(q)
        tdVal := globalTD.Quantile(q)
        writer.WriteString(fmt.Sprintf("Квантиль %.2f: TDigest = %.4f, Настоящие = %.4f, Разница = %.4f\n",
            q, tdVal, realVal, diffs[i]))
    }

    writer.Flush()
}