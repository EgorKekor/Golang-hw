package main

import (
	"runtime"
	"sort"
	"strconv"
	"sync"
	"time"
)

type EOD struct {

}

type CrcData struct {
	data	string
	ind 	int
	afterMd5	bool
}

func CRC(data string, ind int, outp chan CrcData, chanUsr *sync.WaitGroup, fromMd5 bool) {
	defer chanUsr.Done()
	hash := DataSignerCrc32(data)
	outp<-CrcData{hash, ind, fromMd5}
	return
}

func simpleCRC(data string, ind int, outp chan CrcData) {
	hash := DataSignerCrc32(data)
	outp<-CrcData{hash, ind, false}
	return
}

func MD5(data string, ind int, outp chan CrcData, chanUsr *sync.WaitGroup) {
	defer chanUsr.Done()
	hash := DataSignerMd5(data)
	outp<-CrcData{hash, ind, true}
	return
}


func startMd5Worker(inp, outp chan CrcData, chanUsr *sync.WaitGroup, localWorker *sync.WaitGroup) {
	defer localWorker.Done()
	println("Start MD5 worker")
	ticker := time.Tick(11 * time.Millisecond)

	for data := range inp {
		<-ticker
		go MD5(data.data, data.ind, outp, chanUsr)
	}
	println("Kill MD5 worker")
	return
}


func startCrcWorker(inp, outp chan CrcData, chanUsr *sync.WaitGroup, localWorker *sync.WaitGroup) {
	defer localWorker.Done()
	println("Start CRC worker")
	for data := range inp {
		go CRC(data.data, data.ind, outp, chanUsr, data.afterMd5)
	}
	println("Kill CRC worker")
	return
}


func startConstructorWorker(inp chan CrcData, outp chan interface{}, localWorker *sync.WaitGroup) {
	defer localWorker.Done()
	println("Start Constructor worker")
	concat := make(map[int]string)

	for data := range inp {
		if _, ok := concat[data.ind]; ok {		//Если ключь существует
			val, _ := concat[data.ind]			//Взять значение по нему
			if data.afterMd5 {
				result := val + "~" + data.data
				outp<-result
			} else {
				result := data.data + "~" + val
				outp<-result
			}
		} else {
			concat[data.ind] = data.data
		}
	}
	println("Kill Constructor worker")
	return
}


func startMultiConstructorWorker(inp chan CrcData, outp chan interface{}, localWorker *sync.WaitGroup) {
	defer localWorker.Done()
	println("Start MultiConstructor worker")
	concat := make(map[int]string)
	for th := 0; th < 6; th++ {
		data := <-inp
		concat[data.ind] = data.data
	}

	result := ""
	for th := 0; th < 6; th++ {
		result += concat[th]
	}
	outp<-result

	println("Kill MultiConstructor worker")
	return
}


type By func(str1, str2 *string) bool

func (by By) Sort(str []string) {
	sortCfg := &stringSorter {
		obj: str,
		by:  by,
	}
	sort.Sort(sortCfg)
}


type stringSorter struct {
	obj []string
	by  func(str1, str2 *string) bool
}


func (sorter *stringSorter) Len() int {
	return len(sorter.obj)
}

func (sorter *stringSorter) Swap(i, j int) {
	sorter.obj[i], sorter.obj[j] = sorter.obj[j], sorter.obj[i]
}


func (sorter *stringSorter) Less(i, j int) bool {
	return sorter.by(&sorter.obj[i], &sorter.obj[j])
}

//  ==============================================================

func SingleHash(in, out chan interface{}) {
	CrcOutput := make(chan CrcData, 1)		// Все результаты CRC идут сюда
	CrcInput := make(chan CrcData, 1)		// СRС воркер читает отсюда
	Md5Input := make(chan CrcData, 1)		// MD5 воркер читает отсюда

	channelsUsers := sync.WaitGroup{}		// СRC и MD5 функции используют CrcInput Md5Input которые надо закрыть чтобы убить воркеров
	localWorkers := sync.WaitGroup{}

	localWorkers.Add(3)
	go startMd5Worker(Md5Input, CrcInput, &channelsUsers, &localWorkers)
	go startCrcWorker(CrcInput, CrcOutput, &channelsUsers, &localWorkers)
	go startConstructorWorker(CrcOutput, out, &localWorkers)

	i := 0
	for data := range in {
		dataString := strconv.Itoa(data.(int))
		channelsUsers.Add(3)
		Md5Input <- CrcData{dataString, i, true}
		CrcInput <- CrcData{dataString, i, false}
		i++
		runtime.Gosched()
	}

	channelsUsers.Wait()
	close(Md5Input)
	close(CrcInput)
	for len(CrcOutput) > 0 {runtime.Gosched()}		// Новым данным неоткуда взяться, спим пока конструктор их соберет
	close(CrcOutput)								// Это убьет конструктор
	localWorkers.Wait()
	return
}


func MultiHash(in, out chan interface{}) {
	localWorkers := sync.WaitGroup{}
	multi := make([]chan CrcData, 0)

	dataNum := 0
	for data := range in {
		multi = append(multi, make(chan CrcData, 6))
		dataString := data.(string)

		for th := 0; th < 6; th++ {
			go simpleCRC(strconv.Itoa(th) + dataString, th, multi[dataNum])
		}

		localWorkers.Add(1)
		go startMultiConstructorWorker(multi[dataNum], out, &localWorkers)
		dataNum++
	}
	localWorkers.Wait()
}


func CombineResults(in, out chan interface{}) {
	simple := func(str1, str2 *string) bool {
		return *str1 < *str2
	}

	var sortData []string
	for data := range in {
		sortData = append(sortData, data.(string))
	}

	By(simple).Sort(sortData)

	var result string
	for i, str := range sortData {
		result += str
		if i != len(sortData) - 1 {
			result += "_"
		}
	}

	out<-result
}


func wrapper(in, out chan interface{}, worker job) {
	worker(in, out)
	for len(out) > 0 {runtime.Gosched()}	// Пока следующий воркер дочитает
	close(out)
}


func ExecutePipeline(workers ...job) {
	channels := make([]chan interface{}, 0, len(workers))
	runtime.GOMAXPROCS(4)
	for i := 0; i < len(workers) + 1; i++ {
		channels = append(channels, make(chan interface{}, 10))
	}


	for i, worker := range(workers) {
		go wrapper(channels[i], channels[i + 1], worker)
	}
	runtime.Gosched()
	<-channels[len(channels) - 1]
	time.Sleep(10 * time.Millisecond)
	return

}
