package index

import (
	"sort"
	"strings"
	"sync"
)

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
	inData string, filename string, wg *sync.WaitGroup, sem *chan int) error {
	//words := strings.Split(str, " ")
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
