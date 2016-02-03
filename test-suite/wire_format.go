package main

import (
	"fmt"
	"strings"
)

//Makeindexmessage Generates a message to index this package
func MakeIndexMessage(pkg *Package) string {
	dependenciesNames := []string{}

	for _, dep := range pkg.Dependencies {
		dependenciesNames = append(dependenciesNames, dep.Name)
	}

	namesAsString := strings.Join(dependenciesNames, ",")
	return fmt.Sprintf("INDEX|%s|%s", pkg.Name, namesAsString)
}

//MakeRemoveMessage generates a message to remove a pakcage from the server's index
func MakeRemoveMessage(pkg *Package) string {
	return fmt.Sprintf("REMOVE|%s|", pkg.Name)
}

//MakeQueryMessage generates a message to check if a package is currently indexed
func MakeQueryMessage(pkg *Package) string {
	return fmt.Sprintf("QUERY|%s|", pkg.Name)
}
