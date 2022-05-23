package utils

import (
	"encoding/csv"
	"log"
	"os"
)

// CsvRow CSV 文件中的一行
type CsvRow struct {
	// 行号
	RowNo int
	// 所有列组成的切片
	Columns []string
}

// ReadCsv 接收一个 chan 并把 csv 文件的每一行作为 csvRow 输入到 chan 中
func ReadCsv(filepath string, columns int, jumpHeader bool) <-chan CsvRow {
	// 读取 CSV 文件，创建文件指针
	ch := make(chan CsvRow)
	csvFile, err := os.Open(filepath)
	if err != nil {
		log.Println("csv文件打开失败！")
	}

	// 初始化 CSVReader
	// 并按照参数选择是否跳过第一行表头
	csvReader := csv.NewReader(csvFile)
	csvReader.FieldsPerRecord = columns
	if jumpHeader {
		_, err = csvReader.Read()
		if err != nil {
			log.Println("读取表头出错", err)
		}
	}

	// 异步读取，先行返回 chan
	go func() {
		defer func(opencast *os.File) {
			err := opencast.Close()
			if err != nil {
				log.Print("Close failed")
			}
		}(csvFile)
		// 记录行号
		// 由于数据集可能非常大，这里使用逐行读取
		rowNo := 1
		for {
			row, err := csvReader.Read()
			if err != nil {
				break
			}
			csvRow := CsvRow{RowNo: rowNo, Columns: row}
			rowNo++
			ch <- csvRow
		}
		defer close(ch)
	}()
	return ch
}
