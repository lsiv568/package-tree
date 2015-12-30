package main

import (
	"fmt"
	"regexp"
	"strings"
)

type Package struct {
	Name         string
	Processed    bool
	Dependencies []*Package
}

func (pkg *Package) AddDependency(to *Package) {
	pkg.Dependencies = append(pkg.Dependencies, to)
}

type AllPackages struct {
	Packages []*Package
}

var (
	lineMatcher, _ = regexp.Compile("^\\w+: ?(\\w+ *)*")
)

func (allPackages *AllPackages) Names() []string {
	names := make([]string, len(allPackages.Packages))
	for _, p := range allPackages.Packages {
		names = append(names, p.Name)
	}
	return names
}

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

func MakeUnprocessedPackage(name string) *Package {
	return &Package{
		Name:         name,
		Processed:    false,
		Dependencies: make([]*Package, 0),
	}
}

func ParsePackageFromLine(line string) (*Package, error) {
	if !lineMatcher.MatchString(line) {
		return nil, fmt.Errorf("Invalid line: %s", line)
	}

	sanitisedLine := strings.Replace(strings.Trim(line, " "), "  ", " ", 100)
	tokens := strings.Split(sanitisedLine, " ")

	packageName := strings.TrimRight(tokens[0], ":")

	dependenciesNames := tokens[1:len(tokens)]
	dependencies := make([]*Package, len(dependenciesNames))

	for i, name := range dependenciesNames {
		dependencies[i] = MakeUnprocessedPackage(name)
	}
	return &Package{
		Name:         packageName,
		Processed:    true,
		Dependencies: dependencies,
	}, nil
}

func GetPackages() AllPackages {
	return AllPackages{
		Packages: make([]*Package, 0),
	}
}
