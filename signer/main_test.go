package main

import (
	"crypto/md5"
	"fmt"
	"hash/crc32"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

/*
это тест на проверку того что у нас это действительно конвейер
неправильное поведение: накапливать результаты выполнения одной функции, а потом слать их в следующую.
	это не похволяет запускать на конвейере бесконечные задачи
правильное поведение: обеспечить беспрепятственный поток
*/
func TestPipeline(t *testing.T) {

	var ok = true
	var recieved uint32
	freeFlowJobs := []job{
		job(func(in, out chan interface{}) {
			out <- 1
			time.Sleep(10 * time.Millisecond)
			currRecieved := atomic.LoadUint32(&recieved)
			// в чем тут суть
			// если вы накапливаете значения, то пока вся функция не отрабоатет - дальше они не пойдут
			// тут я проверяю, что счетчик увеличился в следующей функции
			// это значит что туда дошло значение прежде чем текущая функция отработала
			if currRecieved == 0 {
				ok = false
			}
		}),
		job(func(in, out chan interface{}) {
			for _ = range in {
				atomic.AddUint32(&recieved, 1)
			}
		}),
	}
	ExecutePipeline(freeFlowJobs...)
	if !ok || recieved == 0 {
		t.Errorf("no value free flow - dont collect them")
	}
}



func TestSigner(t *testing.T) {

	testExpected := "1173136728138862632818075107442090076184424490584241521304_1696913515191343735512658979631549563179965036907783101867_27225454331033649287118297354036464389062965355426795162684_29568666068035183841425683795340791879727309630931025356555_3994492081516972096677631278379039212655368881548151736_4958044192186797981418233587017209679042592862002427381542_4958044192186797981418233587017209679042592862002427381542"
	testResult := "NOT_SET"

	// это небольшая защита от попыток не вызывать мои функции расчета
	// я преопределяю фукции на свои которые инкрементят локальный счетчик
	// переопределение возможо потому что я объявил функцию как переменную, в которой лежит функция
	var (
		DataSignerSalt         string = "" // на сервере будет другое значение
		OverheatLockCounter    uint32
		OverheatUnlockCounter  uint32
		DataSignerMd5Counter   uint32
		DataSignerCrc32Counter uint32
	)
	OverheatLock = func() {
		atomic.AddUint32(&OverheatLockCounter, 1)
		for {
			if swapped := atomic.CompareAndSwapUint32(&dataSignerOverheat, 0, 1); !swapped {
				fmt.Println("OverheatLock happend")
				time.Sleep(time.Second)
			} else {
				break
			}
		}
	}
	OverheatUnlock = func() {
		atomic.AddUint32(&OverheatUnlockCounter, 1)
		for {
			if swapped := atomic.CompareAndSwapUint32(&dataSignerOverheat, 1, 0); !swapped {
				fmt.Println("OverheatUnlock happend")
				time.Sleep(time.Second)
			} else {
				break
			}
		}
	}
	DataSignerMd5 = func(data string) string {
		atomic.AddUint32(&DataSignerMd5Counter, 1)
		OverheatLock()
		defer OverheatUnlock()
		data += DataSignerSalt
		dataHash := fmt.Sprintf("%x", md5.Sum([]byte(data)))
		time.Sleep(10 * time.Millisecond)
		return dataHash
	}
	DataSignerCrc32 = func(data string) string {
		atomic.AddUint32(&DataSignerCrc32Counter, 1)
		data += DataSignerSalt
		crcH := crc32.ChecksumIEEE([]byte(data))
		dataHash := strconv.FormatUint(uint64(crcH), 10)
		time.Sleep(time.Second)
		return dataHash
	}

	inputData := []int{0, 1, 1, 2, 3, 5, 8}
	//inputData := []int{0,1}

	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			for _, fibNum := range inputData {
				out <- fibNum
			}
		}),
		job(func(in, out chan interface{}) {
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
		}),

		//  ======================================================

		job(func(in, out chan interface{}) {
			finish := make(chan bool, 6)
			dataNum := 0
			multi := make([]chan CrcData, 0)

			for data := range in {
				println("input-----", data.(string))

				multi = append(multi, make(chan CrcData, 6))
				dataString := data.(string)

				for th := 0; th < 6; th++ {
					go simpleCRC(strconv.Itoa(th) + dataString, th, multi[dataNum])
				}
				go startMultiConstructorWorker(multi[dataNum], out, finish)
				dataNum++
				runtime.Gosched()
			}
			println("BREAK WORKER 2")

			println("Wait constructors begin")
			for constructorNum := 0; constructorNum < dataNum; {
				<-finish
				constructorNum++
				println(constructorNum)
			}
			println("Constructors end")
		}),

		//  ======================================================

		job(func(in, out chan interface{}) {
			simple := func(str1, str2 *string) bool {
				return *str1 < *str2
			}

			var sortData []string

			for data := range in {
				if _, ok := data.(EOD); ok {
					break
				}
				sortData = append(sortData, data.(string))
			}
			println("BREAK WORKER 3")

			By(simple).Sort(sortData)

			result := ""
			for i, str := range sortData {
				result += str
				if i != len(sortData) - 1 {
					result += "_"
				}
			}

			out<-result
		}),

		//  ======================================================

		job(func(in, out chan interface{}) {
			dataRaw := <-in
			data, ok := dataRaw.(string)
			if !ok {
				t.Error("cant convert result data to string")
			}
			testResult = data
			println(testResult)
		}),
	}

	start := time.Now()

	ExecutePipeline(hashSignJobs...)

	end := time.Since(start)

	expectedTime := 3 * time.Second

	if testExpected != testResult {
		t.Errorf("results not match\nGot: %v\nExpected: %v", testResult, testExpected)
	}

	if end > expectedTime {
		t.Errorf("execition too long\nGot: %s\nExpected: <%s", end, time.Second*3)
	}

	// 8 потому что 2 в SingleHash и 6 в MultiHash
	if int(OverheatLockCounter) != len(inputData) ||
		int(OverheatUnlockCounter) != len(inputData) ||
		int(DataSignerMd5Counter) != len(inputData) ||
		int(DataSignerCrc32Counter) != len(inputData)*8 {
		t.Errorf("not enough hash-func calls")
	}

}


