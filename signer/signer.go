package main

import (
	"runtime"
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

func CRC(data string, ind int, outp chan CrcData, inProc *int, mut sync.Mutex, fromMd5 bool) {		//CRC COUNTER
	println("Start CRC[" + strconv.Itoa(ind) + "]")
	hash := DataSignerCrc32(data)
	outp<-CrcData{hash, ind, fromMd5}
	mut.Lock()
	*inProc--
	mut.Unlock()
	println("Finish CRC[" + strconv.Itoa(ind) + "]")
	return
}

func MD5(data string, ind int, outp chan CrcData, inProc *int, mut sync.Mutex) {		//MD5 COUNTER
	println("Start MD5[" + strconv.Itoa(ind) + "]")
	hash := DataSignerMd5(data)
	outp<-CrcData{hash, ind, true}
	mut.Lock()
	*inProc--
	mut.Unlock()
	println("Finish MD5[" + strconv.Itoa(ind) + "]")
	return
}


func startMd5Worker(inp, outp chan CrcData, inProc *int, mut sync.Mutex) {
	println("Start MD5 worker")
	ticker := time.Tick(11 * time.Millisecond)

	for data := range inp {
		<-ticker
		go MD5(data.data, data.ind, outp, inProc, mut)
	}
	println("Kill MD5 worker")
	return
}


func startCrcWorker(inp, outp chan CrcData, inProc *int, mut sync.Mutex) {
	println("Start CRC worker")
	for data := range inp {
		go CRC(data.data, data.ind, outp, inProc, mut, data.afterMd5)
	}
	println("Kill CRC worker")
	return
}


func startConstructorWorker(inp chan CrcData, outp chan interface{}) {
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
	outp <- EOD{}
	return
}


func startMultiConstructorWorker(inp chan CrcData, outp chan interface{}) {
	println("Start MultiConstructor worker")


	println("Kill MultiConstructor worker")
	return
}





func ExecutePipeline(workers ...job) {
	channels := make([]chan interface{}, 0, len(workers))
	runtime.GOMAXPROCS(4)
	for i := 0; i < len(workers) + 1; i++ {
		channels = append(channels, make(chan interface{}, 10))
	}


	for i, worker := range(workers) {
		go worker(channels[i], channels[i + 1])
	}
	runtime.Gosched()
	time.Sleep(100 * time.Second)


}
