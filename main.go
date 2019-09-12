package main

import (
	"flag"
	"log"
	"os"
)

func main() {

	flagF := flag.Bool("f", false, "Ignore register")
	flagU := flag.Bool("u", false, "Only first")
	flagR := flag.Bool("r", false, "Sort low")
	flagO := flag.Bool("0", false, "Write file")
	flagN := flag.Bool("n", false, "Numbers sort")
	flag.Parse()

	if (*flagF && *flagN) {
		return
	}

	if *flagF {
		println(*flagF)
	}

	fileName :=  flag.Args();
	println("name:", fileName[0], "\n")


	f, err := os.OpenFile(fileName[0], os.O_RDWR | os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}


	//message := strings.NewReader("ReaderString\n")
	//fmt.Printf("%T %s\n", message, message);
	//io.Copy(f, message);

	message2 := "SimpleString\n";
	f.WriteString(message2);

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
