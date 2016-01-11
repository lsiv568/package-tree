package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Serialise converts the payload into something to be sent over the wire to server
func Serialise(pkg *Package) string {
	dependenciesNames := []string{}

	for _, dep := range pkg.Dependencies {
		dependenciesNames = append(dependenciesNames, dep.Name)
	}

	namesAsString := strings.Join(dependenciesNames, ",")
	return fmt.Sprintf("INSTALL|%s|%s", pkg.Name, namesAsString)
}

// Deserialise gets the response from the server and interprets it as success
// or failure
func Deserialise(responseMsg string) (bool, error) {
	result, err := strconv.Atoi(strings.TrimRight(responseMsg, "\n"))
	if err != nil {
		return false, err
	}

	return 0 == result, nil
}
