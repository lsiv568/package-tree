package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
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

type PackageIndexerClient interface {
	Name() string
	Close() error
	Send(msg string) (ResponseCode, error)
}

// PackageIndexerClient connects to the running server.
type TcpPackageIndexerClient struct {
	name string
	conn net.Conn
}

//Close closes the connection to the server.
func (client *TcpPackageIndexerClient) Name() string {
	return client.name
}

//Close closes the connection to the server.
func (client *TcpPackageIndexerClient) Close() error {
	log.Printf("%s disconnecting", client.Name())
	return client.conn.Close()
}

//Send sends amessage to the server using its line-oriented protocol
func (client *TcpPackageIndexerClient) Send(msg string) (ResponseCode, error) {
	log.Printf("%s sending message [%s]", client.Name(), msg)
	_, err := fmt.Fprintln(client.conn, msg)
	log.Printf("%s received [%v]", client.Name(), err)

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
func MakeTcpPackageIndexClient(name string, port int) (PackageIndexerClient, error) {
	host := fmt.Sprintf("localhost:%d", port)
	log.Printf("%s connecting to [%s]", name, host)
	conn, err := net.DialTimeout("tcp", host, time.Duration(1)*time.Second)

	if err != nil {
		return nil, fmt.Errorf("Failed to open connection to [%s]: %#v", host, err)
	}

	return &TcpPackageIndexerClient{
		name: name,
		conn: conn,
	}, nil
}
