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

	for _, pkg := range allPackages.Packages {
		log.Printf("Processing package [%s]", pkg.Name)

		pkgAlreadyInstalled, err := send(conn, Serialise("QUERY", pkg))

		if err != nil {
			test.Failf("When reading %v", err)
		}

		if pkgAlreadyInstalled {
			test.Failf("Pacakge %v was already present", pkg.Name)
		}

		successful, err := send(conn, Serialise("INSTALL", pkg))

		if err != nil {
			test.Failf("When reading %v", err)
		}

		if successful {
			fmt.Println("YAY")
		}

		pkgAlreadyInstalled, err = send(conn, Serialise("QUERY", pkg))

		if err != nil {
			test.Failf("When reading %v", err)
		}

		if !pkgAlreadyInstalled {
			test.Failf("Pacakge %v was not installed", pkg.Name)
		}
	}

}

func send(conn net.Conn, msg string) (bool, error) {
	_, err := fmt.Fprintln(conn, msg)

	if err != nil {
		test.Failf("Error sending message to the server %v", err)
	}

	responseMsg, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		test.Failf("Error reading message from server %v", err)
	}

	wasSuccesful, err := Deserialise(responseMsg)

	if err != nil {
		test.Failf("Error parsing message from server [%s] %v", responseMsg, err)
	}

	return wasSuccesful, nil
}
