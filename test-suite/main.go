package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	port := *flag.Int("port", 8080, "The port your server exposes to clients")
	//	concurrencyLevel := flag.Int("concurrency", 10, "A positive value indicating how many concurrent clients to use")
	flag.Parse()

	test := &TestRun{
		ServerPort: port,
	}

	test.Start()

	test.Phase1()
	test.Phase2()
	//	test.Phase3()

	test.Finish()
}
