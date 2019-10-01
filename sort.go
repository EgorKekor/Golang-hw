package main

import (
	"bufio"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type sortObject struct {
	lines      []string
	sortSelect []string
	column 		int
}

func createSortObject(r io.Reader, columnSort int) (sortObject) {
	obj := sortObject {}
	obj.column = columnSort
	fileScanner := bufio.NewScanner(r)

	for i := 0; fileScanner.Scan(); i++ {
		obj.lines = append(obj.lines, strings.TrimLeft(fileScanner.Text(), " "))
		if columnSort > 0 {
			strScanner := bufio.NewScanner(strings.NewReader(fileScanner.Text()))
			strScanner.Split(bufio.ScanWords)
			wordNum := 1
			for strScanner.Scan() && wordNum < columnSort {
				wordNum++
			}

			if wordNum == columnSort {
				obj.sortSelect = append(obj.sortSelect, strScanner.Text())
			} else {
				obj.sortSelect = append(obj.sortSelect, "")
			}
		} else {
			obj.sortSelect = append(obj.sortSelect, obj.lines[i])
		}
	}

	return obj
}

func (sObj *sortObject) setUniqueMode() {
	var dublicateSelect []string
	var dublicateLines []string

	dublicate := make(map[string]bool)
	for i, orig := range(sObj.sortSelect) {
		if _, ok := dublicate[orig]; !ok {
			dublicate[orig] = true
			dublicateSelect = append(dublicateSelect, sObj.sortSelect[i])
			dublicateLines = append(dublicateLines, sObj.lines[i])
		}
	}
	sObj.sortSelect = dublicateSelect
	sObj.lines = dublicateLines
}

func (sObj *sortObject) setLowerCaseMode() {
	for i := 0; i < len(sObj.lines); i++ {
		sObj.sortSelect[i] = strings.ToLower(sObj.sortSelect[i])
	}
}

func (sObj *sortObject) setNumericMode() {
	var numericLines []string
	var numericSelect []string
	for i, _ := range sObj.sortSelect {
		if num, err := strconv.ParseInt(sObj.sortSelect[i], 10, 64); err == nil {
			numericLines = append(numericLines, sObj.lines[i])
			numericSelect = append(numericSelect, strconv.FormatInt(num, 10))
		} else {
			if err == strconv.ErrRange {
				numericLines = append(numericLines, sObj.lines[i])
				numericSelect = append(numericSelect, strconv.FormatInt(math.MaxInt64, 10))
			}
		}
	}
	sObj.lines = numericLines
	sObj.sortSelect = numericSelect
}

func (sObj *sortObject) writeInFile(ok bool, filename string) {

	var output io.Writer

	if ok {
		resultFile, err := os.OpenFile(filename, os.O_RDWR | os.O_CREATE, 0755)
		if err != nil {
			log.Panicln(err)
			return
		}
		defer func() {
			if err := resultFile.Close(); err != nil {
				log.Panicln(err)
			}
		}()
		output = resultFile
	} else {
		output = os.Stdout
	}


	for _, str := range (sObj.lines) {
		if _, err := io.Copy(output, strings.NewReader(str+"\n")); err != nil {
			log.Panicln(err)
			return
		}
	}
}

//  ===================================================


type By func(str1, str2 *string) bool

func (by By) Sort(srt_ sortObject) {
	sortCfg := &stringSorter {
		obj: srt_,
		by:  by,
	}
	sort.Stable(sortCfg)
}


type stringSorter struct {
	obj sortObject
	by  func(str1, str2 *string) bool
}


func (sorter *stringSorter) Len() int {
	return len(sorter.obj.lines)
}

func (sorter *stringSorter) Swap(i, j int) {
	sorter.obj.lines[i], sorter.obj.lines[j] = sorter.obj.lines[j], sorter.obj.lines[i]
	sorter.obj.sortSelect[i], sorter.obj.sortSelect[j] = sorter.obj.sortSelect[j], sorter.obj.sortSelect[i]
}


func (sorter *stringSorter) Less(i, j int) bool {
	return sorter.by(&sorter.obj.sortSelect[i], &sorter.obj.sortSelect[j])
}

//  =========================================


//func main() {
//
//	flagF := flag.Bool("f", false, "Ignore register")
//	flagU := flag.Bool("u", false, "Unique")
//	flagR := flag.Bool("r", false, "Sort low")
//	flagO := flag.Bool("o", false, "Write file")
//	flagN := flag.Bool("n", false, "Numbers sort")
//	flagK := flag.Int("k", 0, "Col number")
//	flag.Parse()
//
//
//	fileName :=  flag.Args()
//	sourceFile, err := os.OpenFile(fileName[0], os.O_RDONLY, 0755)
//	if err != nil {
//		log.Panicln(err)
//		return
//	}
//	defer func() {
//		if err := sourceFile.Close(); err != nil {
//			log.Panicln(err)
//		}
//	}()
//
//
//	sortingObj := createSortObject(sourceFile, *flagK)
//
//	simple := func(str1, str2 *string) bool {
//		return *str1 < *str2
//	}
//
//	reverse := func(str1, str2 *string) bool {
//		return *str1 > *str2
//	}
//
//
//
//	if *flagF {
//		sortingObj.setLowerCaseMode()
//	}
//
//	if *flagU {
//		sortingObj.setUniqueMode()
//	}
//
//	if *flagN {
//		sortingObj.setNumericMode()
//	}
//
//
//	if *flagR {
//		By(reverse).Sort(sortingObj)
//	} else {
//		By(simple).Sort(sortingObj)
//	}
//
//
//	sortingObj.writeInFile(*flagO, "result.dat")
//
//}
