package main

import (
	"fmt"
	"net"
	"testing"
)

func respondWith(t *testing.T, server net.Listener, responseCode string) {
	for {
		conn, err := server.Accept()
		if err != nil {
			t.Fatalf("Error reading socket: %v", err)
		}
		fmt.Fprintln(conn, responseCode)
	}
}

func TestSend(t *testing.T) {
	goodPort := 8088
	goodServer, err := net.Listen("tcp", fmt.Sprintf(":%d", goodPort))
	defer goodServer.Close()

	if err != nil {
		t.Fatalf("Error opening test server: %v", err)
	}

	go respondWith(t, goodServer, "OK")

	client, err := MakePackageIndexClient(goodPort)
	if err != nil {
		t.Fatalf("Error connecting to server: %v", err)
	}

	responseCode, err := client.Send("A")

	if err != nil {
		t.Errorf("Error sending message to server: %v", err)
	}

	if responseCode == FAIL {
		t.Errorf("Expected responseCode to be 1, got %v", responseCode)
	}

	badPort := 8090
	badServer, err := net.Listen("tcp", fmt.Sprintf(":%d", badPort))
	defer badServer.Close()

	if err != nil {
		t.Fatalf("Error opening test server: %v", err)
	}

	go respondWith(t, badServer, "banana")

	client, err = MakePackageIndexClient(badPort)
	if err != nil {
		t.Fatalf("Error connecting to server: %v", err)
	}

	responseCode, err = client.Send("B")

	if err == nil {
		t.Errorf("No error returned for bad responseCode from server: %#v", responseCode)
	}

}
