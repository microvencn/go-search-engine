package keywords

import (
	"GoSearchEngine/utils"
	"bufio"
	"log"
	"os"
)

var File *os.File = nil

func InitKeyWordsFile() {
	if File == nil {
		fileName := utils.GetPath("/database/words.txt")
		var err error
		File, err = os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
		if err != nil {
			log.Println(err)
		}
	}
}

func GetKeyWordsFile() *os.File {
	if File == nil {
		InitKeyWordsFile()
	}
	return File
}

func CloseKeywordsFile() {
	err := File.Close()
	if err != nil {
		log.Println("关闭 words.txt 失败: ", err)
		return
	}
}

func AddKeyWords(word string) {
	writer := bufio.NewWriter(File)
	_, err := writer.WriteString(word + "\n")
	if err != nil {
		log.Printf("写入关键词 '%s' 失败: %s\n", word, err)
		return
	}
	err = writer.Flush()
	if err != nil {
		log.Println("flush failed: ", err)
		return
	}
}

// GetAllKeyWords 获取所有关键词
func GetAllKeyWords() <-chan string {
	return utils.ReadLineFile(GetKeyWordsFile())
}
