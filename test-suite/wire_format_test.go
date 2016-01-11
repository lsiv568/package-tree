package main

import "testing"

func TestSerialise(t *testing.T) {
	allPackages := AllPackages{}

	aPackage := allPackages.Named("a")

	expectedLine := "INSTALL|a|"
	actualLine := Serialise(aPackage)

	if actualLine != expectedLine {
		t.Errorf("Expected %#v to serialise to [%s], got [%s]", aPackage, expectedLine, actualLine)
	}
}

func TestDeserialise(t *testing.T) {
	message := "0\n"
	expectedResponse := true
	actualResponse, err := Deserialise(message)

	if err != nil {
		t.Fatalf("Error: %#v", err)
	}

	if actualResponse != expectedResponse {
		t.Errorf("Expected [%s]->[%t] , got [%t]", message, actualResponse, expectedResponse)
	}

	message = "false\n"
	actualResponse, err = Deserialise(message)

	if err == nil {
		t.Fatalf("Expected error parsing, [%s], got [%t]", message, actualResponse)
	}
}
