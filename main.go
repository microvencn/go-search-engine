package main

import (
	"GoSearchEngine/fenci"
	"GoSearchEngine/searcher"
	"bufio"
	"fmt"
	"os"
)

func main() {
	// 读取字典
	fenci.ReadDict()
	// 此句进行分词，分词结果和文档存储在 database 文件夹中
	// 若需要重新分词请先清空 database 文件夹
	// 然后运行 index 包下的 index_test

	// 循环接收输入，使用 Ctrl+C 中断即可
	for {
		fmt.Print("Please input: ")
		inputReader := bufio.NewReader(os.Stdin)
		input, err := inputReader.ReadString('\r')
		if err != nil {
			return
		}
		bytes := []byte(input)
		input = string(bytes[0 : len(bytes)-1])
		results := searcher.Simple(input, 0, 10)
		for i := 0; i < len(results); i++ {
			fmt.Println(results[i].IdScore, results[i].Doc)
		}
	}
}
