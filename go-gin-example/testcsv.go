package main

import (
	"os"
	"encoding/csv"
	"go-gin-example/pkg/export"
)

func main() {
	f, err := os.Create(export.GetExcelFullPath + "test.csv111")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// 标识文件的编码格式
	f.WriteString("\xEF\xBB\xBF")  // xEF\xBB\xBF 是 UTF-8 BOM 的 16 进制格式

	w := csv.NewWriter(f)

	data := [][]string{
		{"1", "test1", "test1-1"},
		{"2", "test2", "test2-1"},
		{"3", "test3", "test3-1"},
	}

	w.WriteAll(data)
}
