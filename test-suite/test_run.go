package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

// TestRun controls the current state of the test program.
type TestRun struct {
	ServerPort       int
	StartedAt        time.Time
	ConcurrencyLevel int
	Unluckiness      int
	waiting          sync.WaitGroup
}

// Start starts the test
func (t *TestRun) Start() {
	log.Println("================")
	log.Println(" Starting test ")
	log.Println("================")
	log.Printf("expected server port [%d]", t.ServerPort)
	log.Printf("concurrency level    [%d]", t.ConcurrencyLevel)
	log.Printf("unluckiness          [%d]", t.Unluckiness)
	t.StartedAt = time.Now()
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

//Failf fails the test with a formatted message
func (t *TestRun) Failf(format string, a ...interface{}) {
	t.Fail(fmt.Sprintf(format, a...))
}

//Phase1 will test the server for correctness using a single connection
func (t *TestRun) Phase1() {
	startedAt := time.Now()

	log.Println("TESTRUN Phase1 - Make simple checks for correctness using a single client")
	allPackages := &AllPackages{}

	var packagesWithDependencies []*Package
	for _, pkg := range allPackages.Packages {
		if len(pkg.Dependencies) > 0 {
			packagesWithDependencies = append(packagesWithDependencies, pkg)
		}
	}

	duration := time.Since(startedAt)
	log.Printf("TESTRUN Phase1 - FINISHED (took %dms)", durationInMillis(duration))
}

//Phase2 will index all packages and then remove them using multiple concurrent clients
func (t *TestRun) Phase2() {
	startedAt := time.Now()

	log.Println("TESTRUN Phase2 - brute-force indexes and removes a lot of packages using a single client")

	homebrewPackages, err := BrewToPackages(&AllPackages{})
	if err != nil {
		panic(fmt.Sprintf("Error parsing packages"))
	}

	segmentedPackages := SegmentListPackages(homebrewPackages.Packages, t.ConcurrencyLevel)

	clientCounter := 0
	concurrentBruteforceRemovesAllPackages(clientCounter, t, segmentedPackages)

	clientCounter = clientCounter + t.ConcurrencyLevel
	concurrentBruteforceIndexesPackages(clientCounter, t, segmentedPackages)

	clientCounter = clientCounter + t.ConcurrencyLevel
	concurrentverifyAllPackages(clientCounter, t, segmentedPackages, OK)

	clientCounter = clientCounter + t.ConcurrencyLevel
	concurrentBruteforceRemovesAllPackages(clientCounter, t, segmentedPackages)

	clientCounter = clientCounter + t.ConcurrencyLevel
	concurrentverifyAllPackages(clientCounter, t, segmentedPackages, FAIL)

	duration := time.Since(startedAt)
	log.Printf("TESTRUN Phase2 - FINISHED (took %dms %v)", durationInMillis(duration), duration)
}

//MakeTestRun returns a new instance of a test run.
func MakeTestRun(serverPort int, concurrencyLevel int, unluckiness int) *TestRun {
	return &TestRun{
		ServerPort:       serverPort,
		ConcurrencyLevel: concurrencyLevel,
		Unluckiness:      unluckiness,
	}
}

func bruteforceIndexesPackages(client PackageIndexerClient, packages []*Package, changeOfBeingUnluckyInPercent int) error {
	totalPackages := len(packages)
	log.Printf("%s brute-forcing indexing of %d packages", client.Name(), totalPackages)
	for numPackagesInstalledThisItearion := 0; numPackagesInstalledThisItearion < totalPackages; {
		numPackagesInstalledThisItearion = 0
		for _, pkg := range packages {
			if shouldSomethingBadHappen(changeOfBeingUnluckyInPercent) {
				err := sendBrokenMessage(client)
				if err != nil {
					return err
				}
			}

			msg := MakeIndexMessage(pkg)
			responseCode, err := client.Send(msg)

			if err != nil {
				return fmt.Errorf("%s found error when sending message [%s]: %v", client.Name(), msg, err)
			}

			if responseCode == OK {
				numPackagesInstalledThisItearion = numPackagesInstalledThisItearion + 1
			}
		}
		log.Printf("%s reports %v/%v packages indexed", client.Name(), numPackagesInstalledThisItearion, totalPackages)
	}

	return nil
}

func bruteforceRemovesAllPackages(client PackageIndexerClient, packages []*Package, changeOfBeingUnluckyInPercent int) error {
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

func verifyAllPackages(client PackageIndexerClient, packages []*Package, expectedResponseCode ResponseCode, changeOfBeingUnluckyInPercent int) error {
	totalPackages := len(packages)
	log.Printf("%s querying for %d packages and expecting status code to be [%s]", client.Name(), totalPackages, expectedResponseCode)
	for _, pkg := range packages {
		msg := MakeQueryMessage(pkg)
		responseCode, err := client.Send(msg)
		if err != nil {
			return fmt.Errorf("%s found error when sending message [%s]: %v", client.Name(), msg, err)
		}

		if responseCode != expectedResponseCode {
			return fmt.Errorf("%s expected query for package [%s] to return [%s], got [%s]", client.Name(), pkg.Name, expectedResponseCode, responseCode)
		}
	}

	return nil
}

func makeClient(clientName string, t *TestRun) PackageIndexerClient {
	client, err := MakeTCPPackageIndexClient(clientName, t.ServerPort)
	if err != nil {
		t.Failf("Error opening client to t.ServerPort [%d]: %v", t.ServerPort, err)
	}
	return client
}

func concurrentBruteforceIndexesPackages(clientCounter int, t *TestRun, segmentedPackages [][]*Package) {
	t.waiting.Add(t.ConcurrencyLevel)
	for _, p := range segmentedPackages {
		clientCounter++
		go func(number int, packagesToProcess []*Package) {
			name := fmt.Sprintf("client[%d]", number+1)
			log.Printf("Starting %s", name)
			defer t.waiting.Done()

			client := makeClient(name, t)
			defer client.Close()

			err := bruteforceIndexesPackages(client, packagesToProcess, t.Unluckiness)
			if err != nil {
				t.Failf("%v", err)
			}
		}(clientCounter, p)
	}
	t.waiting.Wait()
}

func concurrentBruteforceRemovesAllPackages(clientCounter int, t *TestRun, segmentedPackages [][]*Package) {
	t.waiting.Add(t.ConcurrencyLevel)
	for _, p := range segmentedPackages {
		clientCounter++
		go func(number int, packagesToProcess []*Package) {
			name := fmt.Sprintf("client[%d]", number+1)
			log.Printf("Starting %s", name)
			defer t.waiting.Done()

			client := makeClient(name, t)
			defer client.Close()

			err := bruteforceRemovesAllPackages(client, packagesToProcess, t.Unluckiness)
			if err != nil {
				t.Failf("%v", err)
			}
		}(clientCounter, p)
	}
	t.waiting.Wait()
}

func concurrentverifyAllPackages(clientCounter int, t *TestRun, segmentedPackages [][]*Package, expectedRepose ResponseCode) {
	t.waiting.Add(t.ConcurrencyLevel)
	for _, p := range segmentedPackages {
		clientCounter++
		go func(number int, packagesToProcess []*Package) {
			name := fmt.Sprintf("client[%d]", number+1)
			log.Printf("Starting %s", name)
			defer t.waiting.Done()

			client := makeClient(name, t)
			defer client.Close()

			err := verifyAllPackages(client, packagesToProcess, expectedRepose, t.Unluckiness)
			if err != nil {
				t.Failf("%v", err)
			}
		}(clientCounter, p)
	}
	t.waiting.Wait()
}

func durationInMillis(d time.Duration) int64 {
	return d.Nanoseconds() / int64(time.Millisecond)
}

func shouldSomethingBadHappen(changeOfBeingUnluckyInPercent int) bool {
	return rand.Intn(100) < changeOfBeingUnluckyInPercent
}

func sendBrokenMessage(client PackageIndexerClient) error {
	msg := MakeBrokenMessage()
	response, err := client.Send(msg)

	if err != nil {
		return fmt.Errorf("%s sent broken message [%s] and expected response code [FAIL], but an error was returned: %v", client.Name(), msg, err)
	}

	if response != FAIL {
		return fmt.Errorf("%s sent broken message [%s] and expected response code [FAIL] but got status code [%s]", client.Name(), msg, response)
	}
	return nil
}
