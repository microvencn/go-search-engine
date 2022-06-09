package utils

import (
	"bufio"
	"io"
	"log"
	"os"
)

func GetPath(path string) string {
	workDir, _ := os.Getwd()
	return workDir + "/src/service" + path
}

func ReadLineFile(file *os.File) <-chan string {
	ch := make(chan string)
	reader := bufio.NewReader(file)
	go func() {
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Println("words.txt close failed")
			}
		}(file)
		defer close(ch)
		for {
			data, _, err := reader.ReadLine()
			if err == io.EOF {
				return
			}
			ch <- string(data)
		}
	}()
	return ch
}
