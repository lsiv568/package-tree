package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

//ResponseCode is the code returned by the sever as a response to our requests
type ResponseCode string

const (
	//OK code
	OK = "OK"

	//FAIL code
	FAIL = "FAIL"

	//ERROR code
	ERROR = "ERROR"

	//UNKNWON code
	UNKNWON = "UNKNWON"
)

// PackageIndexerClient connects to the running server.
type PackageIndexerClient struct {
	conn net.Conn
}

//Close closes the connection to the server.
func (client *PackageIndexerClient) Close() error {
	return client.conn.Close()
}

//Send sends amessage to the server using its line-oriented protocol
func (client *PackageIndexerClient) Send(msg string) (ResponseCode, error) {
	_, err := fmt.Fprintln(client.conn, msg)

	if err != nil {
		return UNKNWON, fmt.Errorf("Error sending message: %v", err)
	}

	responseMsg, err := bufio.NewReader(client.conn).ReadString('\n')

	if err != nil {
		return UNKNWON, fmt.Errorf("Error reading message from server: %v", err)
	}

	returnedString := strings.TrimRight(responseMsg, "\n")

	if returnedString == OK {
		return OK, nil
	}

	if returnedString == FAIL {
		return FAIL, nil
	}

	if returnedString == ERROR {
		return ERROR, nil
	}

	return UNKNWON, fmt.Errorf("Error parsing message from server [%s]: %v", responseMsg, err)
}

// MakePackageIndexClient returns a new instance of the client
func MakePackageIndexClient(port int) (*PackageIndexerClient, error) {
	host := fmt.Sprintf("localhost:%d", port)
	log.Printf("Connecting to [%s]", host)
	conn, err := net.Dial("tcp", host)

	if err != nil {
		return nil, fmt.Errorf("Failed to open connection to [%s]: %#v", host, err)
	}

	return &PackageIndexerClient{
		conn: conn,
	}, nil
}
