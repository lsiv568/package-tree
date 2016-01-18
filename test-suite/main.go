package main

import (
	"flag"
	"fmt"
	"log"
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

	client, err := MakePackageIndexClient(port)
	defer client.Close()

	if err != nil {
		test.Failf("Error opening client to port [%d]: %v", port, err)
	}

	for installedPackages := 0; installedPackages < len(allPackages.Packages); {
		installedPackages = 0
		for _, pkg := range allPackages.Packages {
			responseCode, err := client.Send(MakeQueryMessage(pkg))

			if err != nil {
				test.Failf("When reading %v", err)
			}

			responseCode, err = client.Send(MakeIndexMessage(pkg))

			if err != nil {
				test.Failf("When reading %v", err)
			}

			if responseCode == OK {
				responseCode, err = client.Send(MakeQueryMessage(pkg))
				installedPackages = installedPackages + 1

				if err != nil {
					test.Failf("When reading %v", err)
				}

				if responseCode == FAIL {
					test.Failf("Pacakge %v was not installed", pkg.Name)
				}
			}
		}
		log.Printf("%v packages installed of a total of %v packages", installedPackages, len(allPackages.Packages))
	}

	for installedPackages := len(allPackages.Packages); installedPackages > 0; {
		installedPackages = len(allPackages.Packages)

		for _, pkg := range allPackages.Packages {
			responseCode, err := client.Send(MakeRemoveMessage(pkg))

			if err != nil {
				test.Failf("When reading %v", err)
			}

			if responseCode == OK {
				installedPackages = installedPackages - 1
			}

		}

		log.Printf("%v packages still installed", installedPackages)
	}

	test.Finish()
}
