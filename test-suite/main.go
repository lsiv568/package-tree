package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	port := *flag.Int("port", 8080, "The port your server exposes to clients")
	//	concurrencyLevel := flag.Int("concurrency", 10, "A positive value indicating how many concurrent clients to use")
	flag.Parse()

	allPackages, err := BrewToPackages(&AllPackages{})
	if err != nil {
		panic(fmt.Sprintf("Error parsing packages"))
	}

	test := TestRun{}
	test.Start()

	log.Printf("Connecting to port [%d]", port)
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
	defer conn.Close()

	if err != nil {
		test.Failf("Error connecting to port [%d] (%v)", port, err)
	}

	for _, pkg := range allPackages.Packages {
		log.Printf("Processing package [%s]", pkg.Name)

		msg := Serialise(pkg)

		_, err := fmt.Fprintln(conn, msg)
		if err != nil {
			test.Failf("Error sending packages to the server %v", err)
		}

		responseMsg, err := bufio.NewReader(conn).ReadString('\n')

		if err != nil {
			test.Failf("When reading %v", err)
		}

		successful, err := Deserialise(responseMsg)

		if err != nil {
			test.Failf("When reading %v", err)
		}

		if successful {
			fmt.Println("YAY")
		}

	}

}
