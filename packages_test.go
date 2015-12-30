package main

import (
	"reflect"
	"testing"
)

func TestParsePackageFromLine(t *testing.T) {
	lineWithoutDependencies := "a:"
	expectedPackage := &Package{
		Name:         "a",
		Processed:    true,
		Dependencies: make([]*Package, 0),
	}

	pkg, err := ParsePackageFromLine(lineWithoutDependencies)

	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if !reflect.DeepEqual(pkg, expectedPackage) {
		t.Errorf("Couldn't parse package without dependencies: %v != %v", *pkg, expectedPackage)
	}

	lineWithDependencies := "abcde:  autoconf  automake  cd-discid "
	expectedPackage = &Package{
		Name:      "abcde",
		Processed: true,
		Dependencies: []*Package{
			MakeUnprocessedPackage("autoconf"),
			MakeUnprocessedPackage("automake"),
			MakeUnprocessedPackage("cd-discid"),
		},
	}

	pkg, err = ParsePackageFromLine(lineWithDependencies)

	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if !reflect.DeepEqual(pkg, expectedPackage) {
		t.Errorf("Couldn't parse package with dependencies: %v != %v", *pkg, *expectedPackage)
	}

	brokenLine := "missing tokens"
	_, err = ParsePackageFromLine(brokenLine)

	if err == nil {
		t.Error("Didn't throw error on broken line")
	}
}

func TestAllPackages_Named(t *testing.T) {
	allPackages := AllPackages{}

	aPackage := allPackages.Named("pkg-a")
	theSamePackage := allPackages.Named("pkg-a")
	if aPackage != theSamePackage {
		t.Error("Returning different instances for same package")
	}
}

func TestAddingDependencies(t *testing.T) {
	allPackages := AllPackages{}

	pkg1 := allPackages.Named("pkg-1")
	pkg2 := allPackages.Named("pkg-2")
	pkg3 := allPackages.Named("pkg-3")
	pkg4 := allPackages.Named("pkg-4")

	pkg1.AddDependency(pkg2)
	pkg2.AddDependency(pkg3)
	pkg2.AddDependency(pkg4)
	pkg3.AddDependency(pkg4)

	if !reflect.DeepEqual(pkg1.Dependencies, []*Package{pkg2}) {
		t.Errorf("pkg1 should depend on pkg2")
	}

	if !reflect.DeepEqual(pkg2.Dependencies, []*Package{pkg3, pkg4}) {
		t.Errorf("pkg2 should depend on pkg3 and pkg4")
	}

	if !reflect.DeepEqual(pkg3.Dependencies, []*Package{pkg4}) {
		t.Errorf("pkg3 should depend on pkg4")
	}

	if !reflect.DeepEqual(pkg4.Dependencies, []*Package{}) {
		t.Errorf("pkg4 shouldnt depend on anything")
	}
}
