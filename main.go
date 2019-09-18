package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
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

	for fileScanner.Scan() {
		obj.lines = append(obj.lines, fileScanner.Text())
		if columnSort > 0 {
			strScanner := bufio.NewScanner(strings.NewReader(fileScanner.Text()))
			strScanner.Split(bufio.ScanWords)
			wordNum := 1
			for strScanner.Scan() && wordNum < columnSort {
				wordNum++
			}

			if (wordNum == columnSort) {
				obj.sortSelect = append(obj.sortSelect, strScanner.Text())
			} else {
				obj.sortSelect = append(obj.sortSelect, "")
			}
		}
	}
	if err := fileScanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading file lines:", err)
	}

	if len(obj.sortSelect) == 0 {
		for i := 0; i < len(obj.lines); i++ {
			strScanner := bufio.NewScanner(strings.NewReader(obj.lines[i]))
			for strScanner.Scan() {
				obj.sortSelect = append(obj.sortSelect, strScanner.Text())
			}
		}
	}

	return obj;
}

func (sObj *sortObject) setUniqueMode() {
	var dublicateSelect []string
	var dublicateLines []string
	for orig := 0; orig < len(sObj.lines); orig++ {
		unique := true
		for dub := 0; dub < len(dublicateSelect); dub++ {
			if (sObj.sortSelect[orig] == dublicateSelect[dub]) {
				unique = false
				break;
			}
		}

		if unique {
			dublicateSelect = append(dublicateSelect, sObj.sortSelect[orig])
			dublicateLines = append(dublicateLines, sObj.lines[orig])
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

//  ===================================================


type By func(str1, str2 *string) bool

func (by By) Sort(srt_ sortObject) {
	sortCfg := &stringSorter {
		obj: srt_,
		by:  by,
	}
	sort.Sort(sortCfg)
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


func main() {

	flagF := flag.Bool("f", false, "Ignore register")
	flagU := flag.Bool("u", false, "Only first")
	flagR := flag.Bool("r", false, "Sort low")
	flagO := flag.Bool("o", false, "Write file")
	//flagN := flag.Bool("n", false, "Numbers sort")
	flagK := flag.Int("k", 0, "Col number")
	flag.Parse()


	fileName :=  flag.Args();
	sourceFile, err := os.OpenFile(fileName[0], os.O_RDONLY, 0755)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func() {
		if err := sourceFile.Close(); err != nil {
			log.Fatal(err)
		}
	}()


	sortingObj := createSortObject(sourceFile, *flagK);

	simple := func(str1, str2 *string) bool {
		return *str1 < *str2
	}

	reverse := func(str1, str2 *string) bool {
		return *str1 > *str2
	}



	if *flagF {
		sortingObj.setLowerCaseMode()
	}

	if *flagU {
		sortingObj.setUniqueMode()
	}


	if *flagR {
		By(reverse).Sort(sortingObj)
	} else {
		By(simple).Sort(sortingObj)
	}


	var output io.Writer

	if *flagO {
		// maybe close defer
		resultFile, err := os.OpenFile("result.dat", os.O_RDWR | os.O_CREATE, 0755)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer func() {
			if err := sourceFile.Close(); err != nil {
				log.Fatal(err)
			}
		}()
		output = resultFile
	} else {
		output = os.Stdout
	}


	for i := 0; i < len(sortingObj.lines); i++ {
		io.Copy(output, strings.NewReader(sortingObj.lines[i] + "\n"))
	}

}
