package main

import (
	"fmt"
	"regexp"
	"strings"
)

//Package represents a package and its dependencies
type Package struct {
	Name         string
	Dependencies []*Package
}

//AddDependency makes this package depend on some other
func (pkg *Package) AddDependency(to *Package) {
	pkg.Dependencies = append(pkg.Dependencies, to)
}

//AllPackages is a repository for all known packages
type AllPackages struct {
	//All packages we know of
	Packages []*Package
}

var (
	//Matches well-formed lines from the data file, see data.go
	//and data/brew-dependencies.txt
	lineMatcher, _ = regexp.Compile("^\\w+: ?(\\w+ *)*")
)

// Names returns the names of all known packages
func (allPackages *AllPackages) Names() []string {
	names := make([]string, len(allPackages.Packages))
	for _, p := range allPackages.Packages {
		names = append(names, p.Name)
	}
	return names
}

// Named finds or creates a package with given name. This should be the only
// function used to instantiate packages in production so that we can
// keep a single instance per package.
func (allPackages *AllPackages) Named(name string) *Package {
	var pkg *Package

	for _, p := range allPackages.Packages {
		if p.Name == name {
			pkg = p
		}
	}

	if pkg == nil {
		pkg = MakeUnprocessedPackage(name)
		allPackages.Packages = append(allPackages.Packages, pkg)
	}

	return pkg
}

// MakeUnprocessedPackage is an utility function to
// create a package. Should not be used directly
// from production code, use AllPackages#Named()
func MakeUnprocessedPackage(name string) *Package {
	return &Package{
		Name:         name,
		Dependencies: make([]*Package, 0),
	}
}

// ParseLine pares a single line from the text file, returns the relevant
// tokens as an array. The first element of the array is the package
// name, any subsequent elements are dependencies.
func ParseLine(line string) ([]string, error) {
	if !lineMatcher.MatchString(line) {
		return nil, fmt.Errorf("Invalid line: %s", line)
	}

	sanitisedLine := strings.Replace(strings.Trim(line, " "), "  ", " ", 100)
	tokens := strings.Split(sanitisedLine, " ")

	packageName := strings.TrimRight(tokens[0], ":")

	dependenciesNames := tokens[1:len(tokens)]
	return append([]string{packageName}, dependenciesNames...), nil
}
