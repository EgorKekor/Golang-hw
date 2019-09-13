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

var flagKSortBuffer []string

type By func(ind1, ind2 int) bool
func (by By) Sort(lines[] string) {
	sortCfg := &stringSorter {
		lines: lines,
		by:      by,
	}
	sort.Sort(sortCfg)
}


type stringSorter struct {
	lines[] string
	by      func(ind1, ind2 int) bool
}


func (sorter *stringSorter) Len() int {
	return len(sorter.lines)
}

func (sorter *stringSorter) Swap(i, j int) {
	sorter.lines[i], sorter.lines[j] = sorter.lines[j], sorter.lines[i]
	if len(flagKSortBuffer) > 0 {
		flagKSortBuffer[i], flagKSortBuffer[j] = flagKSortBuffer[j], flagKSortBuffer[i]
	}
}

func (sorter *stringSorter) Less(i, j int) bool {
	return sorter.by(i, j)
}


func main() {

	flagF := flag.Bool("f", false, "Ignore register")
	flagU := flag.Bool("u", false, "Only first")
	flagR := flag.Bool("r", false, "Sort low")
	flagO := flag.Bool("o", false, "Write file")
	flagN := flag.Bool("n", false, "Numbers sort")
	flagK := flag.Int("k", 0, "Col number")
	flag.Parse()
	isFlags := (*flagO || *flagK > 0 || *flagU || *flagR || *flagF || *flagN)


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



	var inputData []string
	fileScanner := bufio.NewScanner(sourceFile)

	for fileScanner.Scan() {
		inputData = append(inputData, fileScanner.Text())
	}
	if err := fileScanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading file lines:", err)
	}



	simple := func(ind1, ind2 int) bool {
		return (inputData[ind1] < inputData[ind2])
	}

	cols := func(ind1, ind2 int) bool {
		return (flagKSortBuffer[ind1] < flagKSortBuffer[ind2])
	}



	if *flagF {
		for i := 0; i < len(inputData); i++ {
			inputData[i] = strings.ToLower(inputData[i])
		}
	}

	if *flagU {
		var dublicate []string
		for orig := 0; orig < len(inputData); orig++ {
			unique := true
			for dub := 0; dub < len(dublicate); dub++ {
				if (inputData[orig] == dublicate[dub]) {
					unique = false
					break;
				}
			}

			if unique {
				dublicate = append(dublicate, inputData[orig])
			}
		}
		inputData = dublicate
	}

	if *flagK > 0 {
		if !(*flagR && *flagN) {
			By(simple).Sort(inputData)
		}

		for i := 0; i < len(inputData); i++ {
			strScanner := bufio.NewScanner(strings.NewReader(inputData[i]))
			strScanner.Split(bufio.ScanWords)
			wordNum := 1
			for strScanner.Scan() && wordNum < *flagK {
				wordNum++
			}

			if (wordNum == *flagK) {
				flagKSortBuffer = append(flagKSortBuffer, strScanner.Text())
			} else {
				flagKSortBuffer = append(flagKSortBuffer, "")
			}
		}
		By(cols).Sort(inputData)
	}


	if !isFlags {
		By(simple).Sort(inputData)
	}


	var output io.Writer

	if *flagO {
		// maybe close defer
		resultFile, _ := os.OpenFile("result.dat", os.O_RDWR | os.O_CREATE, 0755)
		output = resultFile
	} else {
		output = os.Stdout
	}

	for i := 0; i < len(inputData); i++ {
		io.Copy(output, strings.NewReader(inputData[i] + "\n"))
	}
}
