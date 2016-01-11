package main

import (
	"fmt"
	"log"
	"os"
)

var ()

type TestRun struct {
}

func (t *TestRun) Start() {
	log.Println("Stating test...")
}

func (t *TestRun) Fail(reason string) {
	log.Printf("Test failed: %s", reason)
	os.Exit(1)
}

func (t *TestRun) Failf(format string, a ...interface{}) {
	t.Fail(fmt.Sprintln(format, a))
}
