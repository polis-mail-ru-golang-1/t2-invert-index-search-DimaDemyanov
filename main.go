package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"https://github.com/polis-mail-ru-golang-1/t2-invert-index-search-DimaDemyanov/index"
)

type resStruct struct {
	filename string
	Count    int
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("No files found")
		fmt.Println("Use:")
		fmt.Println(os.Args[0] + " [file1 file2 ...]")
		os.Exit(0)
	}

	var fileIndex map[string]Index
	fileIndex = make(map[string]Index, 10)
	for i := 1; i < len(os.Args); i++ {
		err := FileIndexing(fileIndex, os.Args[i])
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	}
	var userStr string
	fmt.Print("Enter your phrase: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	userStr = scanner.Text()
	var resultIdx []resStruct
	words := strings.Split(userStr, " ")
	for i := 0; i < len(words); i++ {
		val, ok := fileIndex[words[i]]
		if ok {
			for j := 0; j < len(val.Files); j++ {
				var isInMap bool = false
				for k := 0; k < len(resultIdx); k++ {
					if resultIdx[k].filename == val.Files[j].Filename {
						resultIdx[k].Count += val.Files[j].Count
						isInMap = true
					}
				}
				if !isInMap {
					tmp := resStruct{val.Files[j].Filename, val.Files[j].Count}
					resultIdx = append(resultIdx, tmp)

				}
			}
		}
	}
	sort.SliceStable(resultIdx, func(i, j int) bool { return resultIdx[i].Count > resultIdx[j].Count })
	for i := 0; i < len(resultIdx); i++ {
		fmt.Println(resultIdx[i].filename, "; совпадений -", resultIdx[i].Count)
	}
}
