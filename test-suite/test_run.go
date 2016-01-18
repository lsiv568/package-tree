package main

import (
	"fmt"
	"log"
	"os"
)

// TestRun controls the current state of the test program.
type TestRun struct {
	ServerPort int
}

// Start starts the test
func (t *TestRun) Start() {
	log.Println("TESTRUN Starting...")
}

// Finish ends the test
func (t *TestRun) Finish() {
	log.Println("================")
	log.Println("All tests passed!")
	log.Println("================")
	log.Println("TESTRUN finished!")
	os.Exit(0)
}

// Fail fails the test
func (t *TestRun) Fail(reason string) {
	log.Printf("TESTRUN failed: %s", reason)
	os.Exit(1)
}

// Failf fails the test
func (t *TestRun) Failf(format string, a ...interface{}) {
	t.Fail(fmt.Sprintf(format, a...))
}

//Phase1 will test the server for correctness using a single connection
func (t *TestRun) Phase1() {
	log.Println("TESTRUN Phase1 - Make simple checks for correctness using a single client")
	allPackages := &AllPackages{}
	client, err := MakePackageIndexClient(t.ServerPort)
	defer client.Close()

	if err != nil {
		t.Failf("Error opening client to t.ServerPort [%d]: %v", t.ServerPort, err)
	}

	var packagesWithDependencies []*Package
	for _, pkg := range allPackages.Packages {
		if len(pkg.Dependencies) > 0 {
			packagesWithDependencies = append(packagesWithDependencies, pkg)
		}
	}

	log.Println("TESTRUN Phase1 - FINISHED")
}

//Phase2 will index all packages and then remove them using a single connection
func (t *TestRun) Phase2() {
	log.Println("TESTRUN Phase2 - Brute-force indexes and removes a lot of packages using a single client")

	client, err := MakePackageIndexClient(t.ServerPort)
	defer client.Close()

	if err != nil {
		t.Failf("Error opening client to t.ServerPort [%d]: %v", t.ServerPort, err)
	}
	homebrewPackages, err := BrewToPackages(&AllPackages{})
	if err != nil {
		panic(fmt.Sprintf("Error parsing packages"))
	}

	t.bruteforceRemovesAllPackages(client, homebrewPackages)
	t.verifyAllPackages(client, homebrewPackages, FAIL)
	t.bruteforceIndexsAllPackages(client, homebrewPackages)
	t.verifyAllPackages(client, homebrewPackages, OK)
	t.bruteforceRemovesAllPackages(client, homebrewPackages)
	t.verifyAllPackages(client, homebrewPackages, FAIL)

	log.Println("TESTRUN Phase2 - FINISHED")
}

func (t *TestRun) bruteforceIndexsAllPackages(client *PackageIndexerClient, packages *AllPackages) {
	totalPackages := len(packages.Packages)
	log.Printf("Brute-forcing indexing of %d packages", totalPackages)
	for installedPackages := 0; installedPackages < totalPackages; {
		installedPackages = 0
		for _, pkg := range packages.Packages {
			responseCode, err := client.Send(MakeQueryMessage(pkg))

			if err != nil {
				t.Failf("When reading %v", err)
			}

			responseCode, err = client.Send(MakeIndexMessage(pkg))

			if err != nil {
				t.Failf("When reading %v", err)
			}

			if responseCode == OK {
				installedPackages = installedPackages + 1
				fmt.Print(".")
			} else {
				fmt.Print("x")
			}
		}
		fmt.Print("\n")
		log.Printf("%v/%v packages indexed", installedPackages, totalPackages)
	}

}

func (t *TestRun) bruteforceRemovesAllPackages(client *PackageIndexerClient, packages *AllPackages) {
	totalPackages := len(packages.Packages)
	log.Printf("Brute-forcing removal of %d packages", totalPackages)
	for installedPackages := totalPackages; installedPackages > 0; {
		installedPackages = len(packages.Packages)

		for _, pkg := range packages.Packages {
			responseCode, err := client.Send(MakeRemoveMessage(pkg))

			if err != nil {
				t.Failf("When reading %v", err)
			}

			if responseCode == OK {
				installedPackages = installedPackages - 1
				fmt.Print(".")
			} else {
				fmt.Print("x")
			}

		}
		fmt.Print("\n")
		log.Printf("%v packages still installed", installedPackages)
	}
}

func (t *TestRun) verifyAllPackages(client *PackageIndexerClient, packages *AllPackages, expectedResponseCode ResponseCode) {
	for _, pkg := range packages.Packages {
		responseCode, err := client.Send(MakeQueryMessage(pkg))

		if err != nil {
			t.Failf("When reading %v", err)
		}

		if responseCode != expectedResponseCode {
			t.Failf("Expected query for package [%s] to return [%s], got [%s]", pkg.Name, expectedResponseCode, responseCode)
		}
	}
}
