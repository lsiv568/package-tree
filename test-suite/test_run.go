package main

import (
	"fmt"
	"log"
	"os"
)

// TestRun controls the current state of the test program.
type TestRun struct {
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
	t.Fail(fmt.Sprintf(format, a))
}
