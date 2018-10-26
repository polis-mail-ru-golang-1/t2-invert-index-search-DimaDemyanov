package index

import (
	"fmt"
	"io/ioutil"
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
	filename string, wg *sync.WaitGroup, sem *chan int) error {
	myBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error occured while reading file:")
		fmt.Println(err)
		return err
	} else {
		str := string(myBytes)
		words := strings.Split(str, " ")
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
				for j := 0; j < len(arrayIndexes[word].Files); j++ {
					if arrayIndexes[word].Files[j].Filename == filename {
						arrayIndexes[word].Files[j].Count++
						sort.SliceStable(arrayIndexes[word].Files, func(i, j int) bool { return arrayIndexes[word].Files[i].Count > arrayIndexes[word].Files[j].Count })
					}
				}

			}
			<-*sem

		}
	}
	wg.Done()
	return nil
}
