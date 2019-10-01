package main

import (
	"runtime"
	"sort"
	"strconv"
	"sync"
)

const operationsAmount = 6

type CrcData struct {
	data	string
	ind 	int
	afterMd5	bool
}

func CRC(inp, outp chan CrcData, calculation *sync.WaitGroup) {
	defer calculation.Done()
	if data, ok := <-inp; ok {
		hash := DataSignerCrc32(data.data)
		outp <- CrcData{hash, data.ind, data.afterMd5}
		return
	}
	return
}

func simpleCRC(data string, ind int, outp chan CrcData) {
	hash := DataSignerCrc32(data)
	outp<-CrcData{hash, ind, false}
	return
}

func MD5(inp, outp chan CrcData, mut *sync.Mutex, calculation *sync.WaitGroup) {
	defer calculation.Done()
	if data, ok := <-inp; ok {
		mut.Lock()
		hash := DataSignerMd5(data.data)
		mut.Unlock()
		outp<-CrcData{hash, data.ind, true}

		return
	}

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
	for th := 0; th < operationsAmount; th++ {
		data := <-inp
		concat[data.ind] = data.data
	}

	result := ""
	for th := 0; th < operationsAmount; th++ {
		result += concat[th]
	}
	outp<-result

	println("Kill MultiConstructor worker")
	return
}

//  ==============================================================

func SingleHash(in, out chan interface{}) {
	CrcOutput := make(chan CrcData, 1)		// Все результаты CRC идут сюда
	CrcInput := make(chan CrcData, 1)		// СRС воркер читает отсюда
	Md5Input := make(chan CrcData, 1)		// MD5 воркер читает отсюда

	md5Mut := sync.Mutex{}

	localWorker := sync.WaitGroup{}
	calculation := sync.WaitGroup{}

	localWorker.Add(1)
	go startConstructorWorker(CrcOutput, out, &localWorker)

	i := 0
	for data := range in {
		dataString := strconv.Itoa(data.(int))
		CrcInput <- CrcData{dataString, i, false}
		Md5Input <- CrcData{dataString, i, false}
		calculation.Add(3)
		go MD5(Md5Input, CrcInput, &md5Mut, &calculation)
		go CRC(CrcInput, CrcOutput, &calculation)
		go CRC(CrcInput, CrcOutput, &calculation)
		i++
		runtime.Gosched()
	}

	calculation.Wait()
	close(Md5Input)
	close(CrcInput)
	close(CrcOutput)
	localWorker.Wait()
	return
}


func MultiHash(in, out chan interface{}) {
	localWorkers := sync.WaitGroup{}
	multi := make([]chan CrcData, 0)

	dataNum := 0
	for data := range in {
		multi = append(multi, make(chan CrcData, 6))
		dataString := data.(string)

		for th := 0; th < operationsAmount; th++ {
			go simpleCRC(strconv.Itoa(th) + dataString, th, multi[dataNum])
		}

		localWorkers.Add(1)
		go startMultiConstructorWorker(multi[dataNum], out, &localWorkers)
		dataNum++
	}
	localWorkers.Wait()
}


func CombineResults(in, out chan interface{}) {
	var sortData []string
	for data := range in {
		sortData = append(sortData, data.(string))
	}

	sort.Strings(sortData)

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
	close(out)
}


func ExecutePipeline(workers ...job) {
	channels := make([]chan interface{}, 0, len(workers))
	runtime.GOMAXPROCS(4)
	for i := 0; i < len(workers) + 1; i++ {
		channels = append(channels, make(chan interface{}, 1))
	}


	for i, worker := range(workers) {
		go wrapper(channels[i], channels[i + 1], worker)
	}
	runtime.Gosched()
	<-channels[len(channels) - 1]
	return

}



//func main() {
//	ch := make(chan int, 10)
//	for i := 0; i < 10; i++ {
//		ch <- i
//	}
//	close(ch)
//	time.Sleep(1000 * time.Millisecond)
//	go func(c chan int){
//		println("in")
//		for i := 0; i < 12; i++ {
//			val, err := <-c
//			println(val, " ", err)
//		}
//	}(ch)
//	runtime.Gosched()
//
//
//	time.Sleep(1000 * time.Millisecond)
//}








