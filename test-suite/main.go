package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

var test = &TestRun{}

func main() {
	log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	port := *flag.Int("port", 8080, "The port your server exposes to clients")
	//	concurrencyLevel := flag.Int("concurrency", 10, "A positive value indicating how many concurrent clients to use")
	flag.Parse()

	allPackages, err := BrewToPackages(&AllPackages{})
	if err != nil {
		panic(fmt.Sprintf("Error parsing packages"))
	}

	test.Start()

	log.Printf("Connecting to port [%d]", port)
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
	defer conn.Close()

	if err != nil {
		test.Failf("Error connecting to port [%d] (%v)", port, err)
	}

	for installedPackages := 0; installedPackages < len(allPackages.Packages); {
		installedPackages = 0

		for _, pkg := range allPackages.Packages {
			result, err := send(conn, Serialise("QUERY", pkg))

			if err != nil {
				test.Failf("When reading %v", err)
			}

			result, err = send(conn, Serialise("INSTALL", pkg))

			if err != nil {
				test.Failf("When reading %v", err)
			}

			if result != 0 {
				result, err = send(conn, Serialise("QUERY", pkg))
				installedPackages = installedPackages + 1

				if err != nil {
					test.Failf("When reading %v", err)
				}

				if result == 0 {
					test.Failf("Pacakge %v was not installed", pkg.Name)
				}
			} else {
				for _, d := range pkg.Dependencies {
					fmt.Printf("Missing: %s -> %#v\n", pkg.Name, d.Name)
				}

			}
		}
		log.Printf("%v packages installed of a total of %v packages", installedPackages, len(allPackages.Packages))

	}
}

func send(conn net.Conn, msg string) (int, error) {
	_, err := fmt.Fprintln(conn, msg)

	if err != nil {
		test.Failf("Error sending message to the server %v", err)
	}

	responseMsg, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		test.Failf("Error reading message from server %v", err)
	}

	result, err := Deserialise(responseMsg)

	if err != nil {
		test.Failf("Error parsing message from server [%s] %v", responseMsg, err)
	}

	return result, nil
}
