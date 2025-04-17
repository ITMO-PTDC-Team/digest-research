package main

import (
    "bufio"
    "fmt"
    "os"
    "path/filepath"
    "strconv"
    "strings"
)

var quantiles = []float64{0.25, 0.5, 0.75, 0.80, 0.90}

// loadTDigestFromFile reads a file with one float per line and feeds values into a new TDigest.
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
    // Указываем папку с входным файлом общего распределения и папку для результатов
    inputFolder := "hour"
    outputFolder := filepath.Join("hour", "result")
    os.MkdirAll(outputFolder, 0755)
    outputFile := filepath.Join(outputFolder, "total.txt")

    // Создаем файл для записи итогового отчёта
    outFile, err := os.Create(outputFile)
    if err != nil {
        panic(err)
    }
    defer outFile.Close()
    writer := bufio.NewWriter(outFile)

    // ЗАГРУЖАЕМ ОБЩЕЕ РАСПРЕДЕЛЕНИЕ ИЗ total.txt
    totalFile := filepath.Join(inputFolder, "total.txt")
    // Создаем "большой" TDigest, добавляя в него все значения из total.txt
    combinedTD := loadTDigestFromFile(totalFile)
    combinedTD.updateCumulative()

    // Загружаем "точное" CDF из того же файла для сравнения (предполагается, что LoadCdfImpl
    // строит кумулятивную функцию распределения с использованием точных значений из отсортированного файла)
    realCdf, err := LoadCdfImpl(totalFile)
    if err != nil {
        panic(err)
    }

    // Сравниваем квантильные оценки TDigest и "точного" распределения
    diffs := CompareCDF(realCdf, combinedTD, quantiles)

    // Записываем в файл информацию о количестве центроид и результаты сравнения квантилей
    writer.WriteString(fmt.Sprintf("Итоговый TDigest: Processed=%d, Unprocessed=%d, Cumulative=%d\n",
        len(combinedTD.processed), len(combinedTD.unprocessed), len(combinedTD.cumulative)))
    writer.WriteString("\nСравнение квантилей (TDigest vs настоящие):\n")
    for i, q := range quantiles {
        realVal := realCdf.Quantile(q)
        tdVal := combinedTD.Quantile(q)
        writer.WriteString(fmt.Sprintf("Квантиль %.2f: TDigest = %.4f, Настоящие = %.4f, Относительная разница = %.4f\n",
            q, tdVal, realVal, diffs[i]))
    }

    writer.Flush()
}