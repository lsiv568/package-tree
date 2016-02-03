package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func main() {
	log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	port := *flag.Int("port", 8080, "The port your server exposes to clients")
	concurrencyLevel := *flag.Int("concurrency", 5, "A positive value indicating how many concurrent clients to use")
	randomSeed := *flag.Int64("seed", 42, "A positive value used to seed the random number generator")
	flag.Parse()

	rand.Seed(randomSeed)

	test := MakeTestRun(port, concurrencyLevel)

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	test.Start()

	test.Phase1()
	test.Phase2()
	//	test.Phase3()

	test.Finish()
}
