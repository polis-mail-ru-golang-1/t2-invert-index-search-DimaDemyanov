package index

import (
	"github.com/polis-mail-ru-golang-1/t2-invert-index-search-DimaDemyanov/filesIn"
	"sort"
	"strings"
	"sync"
)

type ResStruct struct {
	Filename string
	Count    int
}

//Index хранит слово и файлы с весами,содержащие это слово
type Index struct {
	Word  string
	Files []ExtFiles
}

//ExtFiles хранит имя файла и вес этого файла
//(для каждого слова - количество встреч в данном файле)
type ExtFiles struct {
	Filename string
	Count    int
}

//FileIndexing обновляет стркутуру обратного индекса в файле filename
func FileIndexing(arrayIndexes map[string]Index,
	filePath string, filename string,
	wg *sync.WaitGroup, sem *chan int) error {
	//words := strings.Split(str, " ")
	inData, _ := filesIn.ReadData(filePath)
	str := inData
	words := strings.Fields(str)
	for i := 0; i < len(words); i++ {
		*sem <- 1
		word := words[i]
		_, ok := arrayIndexes[word]
		if !ok {
			newWordIdx := Index{Word: word}
			newFile := ExtFiles{filename, 1}
			newWordIdx.Files = append(newWordIdx.Files, newFile)
			arrayIndexes[word] = newWordIdx
		} else {
			isFile := false
			for j := 0; j < len(arrayIndexes[word].Files); j++ {
				if arrayIndexes[word].Files[j].Filename == filename {
					arrayIndexes[word].Files[j].Count++
					isFile = true
					sort.SliceStable(arrayIndexes[word].Files, func(i, j int) bool { return arrayIndexes[word].Files[i].Count > arrayIndexes[word].Files[j].Count })
				}
			}
			if !isFile {
				x := arrayIndexes[word]
				x.Files = append(arrayIndexes[word].Files, ExtFiles{filename, 1})
				arrayIndexes[word] = x
			}

		}
		<-*sem

	}
	wg.Done()
	return nil
}

func PhraseIndexing(userStr string, fileIndex map[string]Index) []ResStruct {
	var resultIdx []ResStruct
	words := strings.Fields(userStr)
	for i := 0; i < len(words); i++ {
		val, ok := fileIndex[words[i]]
		if ok {
			for j := 0; j < len(val.Files); j++ {
				var isInMap bool = false
				for k := 0; k < len(resultIdx); k++ {
					if resultIdx[k].Filename == val.Files[j].Filename {
						resultIdx[k].Count += val.Files[j].Count
						isInMap = true
					}
				}
				if !isInMap {
					tmp := ResStruct{val.Files[j].Filename, val.Files[j].Count}
					resultIdx = append(resultIdx, tmp)

				}
			}
		}
	}

	sort.SliceStable(resultIdx, func(i, j int) bool { return resultIdx[i].Count > resultIdx[j].Count })
	return resultIdx
}
