package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

// TestRun controls the current state of the test program.
type TestRun struct {
	ServerPort       int
	StartedAt        time.Time
	ConcurrencyLevel int
}

// Start starts the test
func (t *TestRun) Start() {
	log.Println("================")
	log.Println(" Starting test ")
	log.Println("================")
	t.StartedAt = time.Now()
	t.ConcurrencyLevel = 1
	log.Println("TESTRUN Starting...")
}

// Finish ends the test
func (t *TestRun) Finish() {
	duration := time.Since(t.StartedAt)
	log.Println("================")
	log.Println("All tests passed!")
	log.Println("================")
	log.Printf("TESTRUN finished! (took %dms)", durationInMillis(duration))
	os.Exit(0)
}

// Fail fails the test
func (t *TestRun) Fail(reason string) {
	duration := time.Since(t.StartedAt)
	log.Println("================")
	log.Println("  Test FAILED!  ")
	log.Println("================")
	log.Printf("Test failed (took %dms)\n%s", durationInMillis(duration), reason)
	os.Exit(0)
	os.Exit(1)
}

// Failf fails the test
func (t *TestRun) Failf(format string, a ...interface{}) {
	t.Fail(fmt.Sprintf(format, a...))
}

//Phase1 will test the server for correctness using a single connection
func (t *TestRun) Phase1() {
	startedAt := time.Now()

	log.Println("TESTRUN Phase1 - Make simple checks for correctness using a single client")
	allPackages := &AllPackages{}
	client, err := MakeTcpPackageIndexClient("-", t.ServerPort)
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

	duration := time.Since(startedAt)
	log.Printf("TESTRUN Phase1 - FINISHED (took %dms)", durationInMillis(duration))
}

//Phase2 will index all packages and then remove them using a single connection
func (t *TestRun) Phase2() {
	startedAt := time.Now()

	log.Println("TESTRUN Phase2 - Brute-force indexes and removes a lot of packages using a single client")

	homebrewPackages, err := BrewToPackages(&AllPackages{})
	if err != nil {
		panic(fmt.Sprintf("Error parsing packages"))
	}

	segmentedPackages := SegmentListPackages(homebrewPackages.Packages, t.ConcurrencyLevel)

	clientCounter := 0

	var waiting sync.WaitGroup

	waiting.Add(t.ConcurrencyLevel)
	for _, p := range segmentedPackages {
		clientCounter++
		go func(number int, packagesToProcess []*Package) {
			name := fmt.Sprintf("client[%d]", number+1)
			log.Printf("Starting %s", name)
			defer waiting.Done()

			client := makeClient(name, t)
			defer client.Close()

			err = BruteforceRemovesAllPackages(client, packagesToProcess)
			if err != nil {
				t.Failf("%v", err)
			}
		}(clientCounter, p)
	}
	waiting.Wait()

	waiting.Add(t.ConcurrencyLevel)
	for _, p := range segmentedPackages {
		clientCounter++
		go func(number int, packagesToProcess []*Package) {
			name := fmt.Sprintf("client[%d]", number+1)
			log.Printf("Starting %s", name)
			defer waiting.Done()

			client := makeClient(name, t)
			defer client.Close()

			err = BruteforceIndexesPackages(client, packagesToProcess)
			if err != nil {
				t.Failf("%v", err)
			}
		}(clientCounter, p)
	}
	waiting.Wait()

	duration := time.Since(startedAt)
	log.Printf("TESTRUN Phase2 - FINISHED (took %dms %v)", durationInMillis(duration), duration)

}

func durationInMillis(d time.Duration) int64 {
	return d.Nanoseconds() / int64(time.Millisecond)
}

func BruteforceIndexesPackages(client PackageIndexerClient, packages []*Package) error {
	totalPackages := len(packages)
	log.Printf("Brute-forcing indexing of %d packages", totalPackages)
	for installedPackages := 0; installedPackages < totalPackages; {
		for _, pkg := range packages {
			responseCode, err := client.Send(MakeIndexMessage(pkg))

			if err != nil {
				return fmt.Errorf("When reading: %v", err)
			}

			if responseCode == OK {
				installedPackages = installedPackages + 1
			}
		}
		log.Printf("%v/%v packages indexed", installedPackages, totalPackages)
	}

	return nil
}

func BruteforceRemovesAllPackages(client PackageIndexerClient, packages []*Package) error {
	totalPackages := len(packages)
	log.Printf("%s brute-forcing removal of %d packages", client.Name(), totalPackages)
	for installedPackages := totalPackages; installedPackages > 0; {
		installedPackages = totalPackages

		for _, pkg := range packages {
			msg := MakeRemoveMessage(pkg)
			responseCode, err := client.Send(msg)
			if err != nil {
				return fmt.Errorf("%s found error when sending message [%s]: %v", client.Name(), msg, err)
			}

			if responseCode == OK {
				installedPackages = installedPackages - 1
			}

		}
		log.Printf("%s reports %d/%d packages still installed", client.Name(), installedPackages, totalPackages)
	}
	return nil
}

func VerifyAllPackages(client PackageIndexerClient, packages []*Package, expectedResponseCode ResponseCode) error {
	totalPackages := len(packages)
	log.Printf("Querying for %d packages and expecting status code to be [%s]", totalPackages, expectedResponseCode)
	for _, pkg := range packages {
		responseCode, err := client.Send(MakeQueryMessage(pkg))

		if err != nil {
			return fmt.Errorf("When reading: %v", err)
		}

		if responseCode != expectedResponseCode {
			return fmt.Errorf("Expected query for package [%s] to return [%s], got [%s]", pkg.Name, expectedResponseCode, responseCode)
		}
	}

	return nil
}

func makeClient(clientName string, t *TestRun) PackageIndexerClient {
	client, err := MakeTcpPackageIndexClient(clientName, t.ServerPort)
	if err != nil {
		t.Failf("Error opening client to t.ServerPort [%d]: %v", t.ServerPort, err)
	}
	return client
}
