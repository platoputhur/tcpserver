package main

import (
	"fmt"
	"net"
	"os"
)

type clientSeverConfig struct {
	clientName string
	conn       net.Conn
}

type clientSever interface {
	readMessage()
	sayBye()
}

func main() {
	//Setting up local port to listen to
	service := ":4564"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkErr(err)

	//Setting up the listener
	listener, err := net.ListenTCP("tcp", tcpAddr)
	fmt.Println("Server started and is currently listening on", service)

	//setting up connection and accept all clients
	var csc clientSeverConfig
	for {
		csc.conn, err = listener.Accept()
		if err != nil {
			continue
		}
		fmt.Printf("Client from %s has been connected!\n", csc.conn.RemoteAddr())
		//opening the message in another go routine
		go csc.readMessage()
	}
}

func (c *clientSeverConfig) readMessage() {
	defer c.sayBye()
	//Making a bufer to read the incoming data
	buf := make([]byte, 1024)
	for {
		n, err := c.conn.Read(buf[0:])
		if err != nil {
			return
		}
		fmt.Print(string(buf[0:n]))
		//Writing confirmation to client terminal
		_, err = c.conn.Write([]byte("Message reached server.\n"))
	}
}

func (c *clientSeverConfig) sayBye() {
	defer c.conn.Close()
	bye := "Bye " + c.conn.RemoteAddr().String()
	fmt.Println(bye)
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
