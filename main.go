package main

import (
	"fmt"
	"github.com/polis-mail-ru-golang-1/t2-invert-index-search-DimaDemyanov/index"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	fileIndex    = make(map[string]index.Index, 10)
	sliceOfNames []string
	portNumber   int
)

func main() {
	arg := os.Args[1:]

	if len(arg) != 2 {
		fmt.Println("Wrong number of arguments.")
		os.Exit(1)
	}
	argPort := strings.TrimLeft(arg[1], ":")
	portNumber, err := strconv.Atoi(argPort)
	if err != nil {
		fmt.Println("Wrong port nubmer: " + "\"" + argPort + "\"")
		os.Exit(2)
	}
	sliceOfNames = getFilesInDir(arg[0])
	var sem = make(chan int, 1)
	var wg sync.WaitGroup
	wg.Add(len(sliceOfNames) - 1)
	for i := 1; i < len(sliceOfNames); i++ {
		go index.FileIndexing(fileIndex, arg[0]+"/"+sliceOfNames[i], sliceOfNames[i], &wg, &sem)
	}
	wg.Wait()
	http.HandleFunc("/", inputForm)
	http.ListenAndServe(":"+strconv.Itoa(portNumber), nil)
}

func inputForm(w http.ResponseWriter, r *http.Request) {
	phrase := r.URL.Query().Get("phrase")

	if phrase != "" {
		fmt.Println(fileIndex)
		resultIdx := index.PhraseIndexing(phrase, fileIndex)
		//fmt.Println(resultIdx)
		for i := 0; i < len(resultIdx); i++ {
			fmt.Fprintln(w, resultIdx[i].Filename, "; совпадений -", resultIdx[i].Count)
		}
	}
}

// выделение имён файлов в заданной директории в слайс
func getFilesInDir(dir string) []string {

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	sliceOfNames := make([]string, 0)
	for _, file := range files {
		if !file.IsDir() {
			sliceOfNames = append(sliceOfNames, file.Name())
		}
	}

	return sliceOfNames
}
