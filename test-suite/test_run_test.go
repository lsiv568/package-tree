package main

import (
	"fmt"
	"testing"
)

type stubClient struct {
	WhatToReturn  ResponseCode
	NumberOfCalls int
	IsCosed       bool
}

//Name returns a hardcoded name
func (client stubClient) Name() string {
	return "stub"
}

//Close does nothing
func (client stubClient) Close() error {
	return nil
}

//Send returns the expected return value and increments the call count
func (client *stubClient) Send(msg string) (ResponseCode, error) {
	client.NumberOfCalls++
	return client.WhatToReturn, nil
}

func TestBruteforceIndexesPackages(t *testing.T) {
	allPackages := &AllPackages{}
	expectedMessages := 20
	for i := 0; i < expectedMessages; i++ {
		allPackages.Named(fmt.Sprintf("pkg-%d", i))
	}

	aStubClient := &stubClient{WhatToReturn: OK}

	BruteforceIndexesPackages(aStubClient, []*Package{})

	if aStubClient.NumberOfCalls != 0 {
		t.Errorf("Expected [%d] calls, got [%d]", expectedMessages, aStubClient.NumberOfCalls)
	}

	aStubClient = &stubClient{WhatToReturn: OK}

	BruteforceIndexesPackages(aStubClient, allPackages.Packages)

	if aStubClient.NumberOfCalls != expectedMessages {
		t.Errorf("Expected [%d] calls, got [%d]", expectedMessages, aStubClient.NumberOfCalls)
	}
}

func TestBruteforceRemovesAllPackages(t *testing.T) {
	allPackages := &AllPackages{}
	expectedMessages := 200
	for i := 0; i < expectedMessages; i++ {
		allPackages.Named(fmt.Sprintf("pkg-%d", i))
	}

	aStubClient := &stubClient{WhatToReturn: OK}

	BruteforceRemovesAllPackages(aStubClient, []*Package{})

	if aStubClient.NumberOfCalls != 0 {
		t.Errorf("Expected [%d] calls, got [%d]", expectedMessages, aStubClient.NumberOfCalls)
	}

	aStubClient = &stubClient{WhatToReturn: OK}

	BruteforceRemovesAllPackages(aStubClient, allPackages.Packages)

	if aStubClient.NumberOfCalls != expectedMessages {
		t.Errorf("Expected [%d] calls, got [%d]", expectedMessages, aStubClient.NumberOfCalls)
	}
}

func TestVerifyAllPackages(t *testing.T) {
	allPackages := &AllPackages{}
	expectedMessages := 200
	for i := 0; i < expectedMessages; i++ {
		allPackages.Named(fmt.Sprintf("pkg-%d", i))
	}

	aStubClient := &stubClient{WhatToReturn: OK}

	VerifyAllPackages(aStubClient, []*Package{}, OK)

	if aStubClient.NumberOfCalls != 0 {
		t.Errorf("Expected [%d] calls, got [%d]", expectedMessages, aStubClient.NumberOfCalls)
	}

	aStubClient = &stubClient{WhatToReturn: OK}

	VerifyAllPackages(aStubClient, allPackages.Packages, FAIL)

	if aStubClient.NumberOfCalls != 1 {
		t.Errorf("Expected to stop after the first failed call, got [%d] calls", aStubClient.NumberOfCalls)
	}
}
