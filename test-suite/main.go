package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
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
	if err != nil {
		test.Fail(fmt.Sprintf("Error connecting to port [%d] (%v)", port, err))
	}

	for _, pkg := range allPackages.Packages {
		dependenciesNames := []string{}

		log.Printf("Processing package [%s]", pkg.Name)

		for _, dep := range pkg.Dependencies {
			dependenciesNames = append(dependenciesNames, dep.Name)
		}

		namesAsString := strings.Join(dependenciesNames, ",")

		msg := fmt.Sprintf("INSTALL|%s|%s", pkg.Name, namesAsString)
		fmt.Fprintln(conn, msg)

		responseMsg, err := bufio.NewReader(conn).ReadString('\n')
		fmt.Println(strings.TrimRight(responseMsg, "\n"))
		returned, err := strconv.Atoi(strings.TrimRight(responseMsg, "\n"))

		if err != nil {
			test.Fail(fmt.Sprintf("When reading %v", err))
		}

		if returned != 0 {
			fmt.Printf("RETURNED %#v", returned)
		}
	}

	conn.Close()
}
