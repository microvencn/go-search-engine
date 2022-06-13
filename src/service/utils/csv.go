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
	Columns    []string
	TotalRowNo int
}

// ReadCsv 接收一个 chan 并把 csv 文件的每一行作为 csvRow 输入到 chan 中
func ReadCsv(columns int, jumpHeader bool, filepath ...string) <-chan CsvRow {
	ch := make(chan CsvRow)

	// 异步读取，先行返回 chan
	go func() {
		totalRowNo := 1
		for _, file := range filepath {
			log.Println("Reading:", file)
			// 读取 CSV 文件，创建文件指针
			csvFile, err := os.Open(file)
			if err != nil {
				log.Println("csv文件打开失败！")
			}

			// 初始化 CSVReader
			// 并按照参数选择是否跳过第一行表头
			csvReader := csv.NewReader(csvFile)
			csvReader.FieldsPerRecord = columns

			var rowNo int
			if jumpHeader {
				rowNo = 1
				//_, err = csvReader.Read()
				//if err != nil {
				//	log.Println("读取表头出错", err)
				//}
			} else {
				rowNo = 0
			}

			rows, _ := csvReader.ReadAll()

			for ; rowNo < len(rows); rowNo++ {
				csvRow := CsvRow{TotalRowNo: totalRowNo, RowNo: rowNo, Columns: rows[rowNo]}
				//if rowNo%100 == 0 {
				//	fmt.Println(rowNo)
				//}
				totalRowNo++
				ch <- csvRow
			}
			// 记录行号
			// 由于数据集可能非常大，这里使用逐行读取
			//rowNo := 1
			//for {
			//	row, err2 := csvReader.Read()
			//	if err2 != nil {
			//		break
			//	}
			//	csvRow := CsvRow{RowNo: rowNo, Columns: row}
			//	rowNo++
			//	if rowNo%100 == 0 {
			//		fmt.Println(rowNo)
			//	}
			//
			//	ch <- csvRow
			//}

			csvFile.Close()
		}
		close(ch)
	}()

	return ch
}
