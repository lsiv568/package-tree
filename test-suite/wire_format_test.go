package main

import "testing"

func TestSerialise(t *testing.T) {
	allPackages := AllPackages{}

	aPackage := allPackages.Named("a")

	action := "INSTALL"
	expectedLine := "INSTALL|a|"
	actualLine := Serialise(action, aPackage)

	if actualLine != expectedLine {
		t.Errorf("Expected %#v to serialise to [%s], got [%s]", aPackage, expectedLine, actualLine)
	}

	action = "QUERY"
	expectedLine = "QUERY|a|"
	actualLine = Serialise(action, aPackage)

	if actualLine != expectedLine {
		t.Errorf("Expected %#v to serialise to [%s], got [%s]", aPackage, expectedLine, actualLine)
	}
}
