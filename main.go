package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/user"
	"path"
	"time"
)

var addr = flag.String("addr", ":1818", "learn to address:port")

func main() {
	flag.Parse()

	usr, err := user.Current()
	if err != nil {
		fmt.Println("Error on getting current user:", err)
	}

	replyFilePath := path.Join(usr.HomeDir, "reply.data")

	listener, err := net.Listen("tcp", *addr)
	if err != nil {
		fmt.Println("Error on initializing listener:", err)
		os.Exit(1)
	}

	defer listener.Close()

	if err != nil {
		fmt.Println("Error on getting hostname:", err)
		os.Exit(1)
	}

	fmt.Println("Listening on - ", *addr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		replyData, err := ioutil.ReadFile(replyFilePath)

		if err != nil {
			fmt.Println("cannot read reply from:", replyFilePath)
		}

		// Handle connections in a new goroutine.
		go handleRequest(conn, string(replyData))
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn, reply string) {
	fmt.Println(
		time.Now().UTC().Format("2006-01-02 15:04:05"),
		": Hanlding connection from - ",
		conn.RemoteAddr())

	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	reqLen, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	} else {
		fmt.Println("get request:", string(buf[:reqLen]))
	}

	conn.Write([]byte(reply))
	fmt.Println("send reply:", reply)
	conn.Close()
}
