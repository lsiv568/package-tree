package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Serialise converts the payload into something to be sent over the wire to server
func Serialise(action string, pkg *Package) string {
	dependenciesNames := []string{}

	for _, dep := range pkg.Dependencies {
		dependenciesNames = append(dependenciesNames, dep.Name)
	}

	namesAsString := strings.Join(dependenciesNames, ",")
	return fmt.Sprintf("%s|%s|%s", action, pkg.Name, namesAsString)
}

// Deserialise gets the response from the server and interprets it as success
// or failure
func Deserialise(responseMsg string) (int, error) {
	result, err := strconv.Atoi(strings.TrimRight(responseMsg, "\n"))
	if err != nil {
		return -1, err
	}

	return result, nil
}
