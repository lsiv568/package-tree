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

func TestDeserialise(t *testing.T) {
	message := "0\n"
	expectedResponse := 0
	actualResponse, err := Deserialise(message)

	if err != nil {
		t.Fatalf("Error: %#v", err)
	}

	if actualResponse != expectedResponse {
		t.Errorf("Expected [%s]->[%v] , got [%v]", message, actualResponse, expectedResponse)
	}

	message = "false\n"
	actualResponse, err = Deserialise(message)

	if err == nil {
		t.Fatalf("Expected error parsing, [%s], got [%v]", message, actualResponse)
	}
}
