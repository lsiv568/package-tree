package main

import (
	"fmt"
	"log"
	"os"
)

var ()

// TestRun controls the current state of the test program.
type TestRun struct {
}

// Start starts the test
func (t *TestRun) Start() {
	log.Println("Stating test...")
}

// Fail fails the test
func (t *TestRun) Fail(reason string) {
	log.Printf("Test failed: %s", reason)
	os.Exit(1)
}

// Failf fails the test
func (t *TestRun) Failf(format string, a ...interface{}) {
	t.Fail(fmt.Sprintf(format, a))
}
