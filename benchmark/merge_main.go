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
    // Папки с исходными данными и для результатов
    inputFolder := "hour"
    outputFolder := filepath.Join("result_hour", "result")
    os.MkdirAll(outputFolder, 0755)
    outputFile := filepath.Join(outputFolder, "merge_1000_123.txt")
    
    // Файл для записи итогового отчёта
    outFile, err := os.Create(outputFile)
    if err != nil {
        panic(err)
    }
    defer outFile.Close()
    writer := bufio.NewWriter(outFile)

    // Создаем глобальный TDigest для объединения всех групп
    globalTD := New()

    // Обработка 60 групп
    for group := 1; group <= 60; group++ {
        groupTD := New()
        // Обработка 60 файлов внутри группы
        for sub := 1; sub <= 60; sub++ {
            filename := filepath.Join(inputFolder, fmt.Sprintf("merge_%d_%d.txt", group, sub))
            td := loadTDigestFromFile(filename)
            groupTD.Merge(td)
            groupTD.updateCumulative()
        }

        // Запись информации по группе
        groupInfo := fmt.Sprintf("Group %d: Processed=%d, Unprocessed=%d, Cumulative=%d\n",
            group, len(groupTD.processed), len(groupTD.unprocessed), len(groupTD.cumulative))
        writer.WriteString(groupInfo)

        // Сливаем групповое распределение в глобальное
        globalTD.Merge(groupTD)
        globalTD.updateCumulative()
        // Также можно записывать промежуточную информацию после слияния группы
        globalInfo := fmt.Sprintf("After merging group %d: Processed=%d, Unprocessed=%d, Cumulative=%d\n",
            group, len(globalTD.processed), len(globalTD.unprocessed), len(globalTD.cumulative))
        writer.WriteString(globalInfo)
    }

    // Загружаем реальное распределение из файла hour/total.txt
    realCdf, err := LoadCdfImpl(filepath.Join(inputFolder, "total.txt"))
    if err != nil {
        panic(err)
    }

    diffs := CompareCDF(realCdf, globalTD, quantiles)

    // Запись сравнения квантилей
    writer.WriteString("\nСравнение квантилей (tdigest vs настоящие):\n")
    for i, q := range quantiles {
        realVal := realCdf.Quantile(q)
        tdVal := globalTD.Quantile(q)
        writer.WriteString(fmt.Sprintf("Квантиль %.2f: tdigest = %.4f, настоящие = %.4f, относительная разница = %.4f\n",
            q, tdVal, realVal, diffs[i]))
    }

    writer.Flush()
}