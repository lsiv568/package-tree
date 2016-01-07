package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
)

var (
	allPackages = AllPackages{}
)

func main() {
	port := flag.Int("port", 8080, "The port your server exposes to clients")
	//	concurrencyLevel := flag.Int("concurrency", 10, "A positive value indicating how many concurrent clients to use")

	log.Printf("Connecting to port [%d]", port)
	conn, err := net.Dial("tcp", "localhost:8080")
	flag.Parse()
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}

	for _, pkg := range allPackages.Packages {
		dependenciesNames := []string{}

		for _, dep := range pkg.Dependencies {
			dependenciesNames = append(dependenciesNames, dep.Name)
		}

		namesAsString := strings.Join(dependenciesNames, ",")

		fmt.Fprintf(conn, "INSTALL|%s|%s\n", pkg.Name, namesAsString)
		status, err := bufio.NewReader(conn).ReadString('\n')

		if err != nil {
			panic(fmt.Sprintf("%v", err))
		}

		fmt.Println(status)
	}
}
