package main

import (
	"fmt"
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
