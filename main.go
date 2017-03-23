package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var addr = flag.String("addr", ":1818", "learn to address:port")

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp", *addr)
	if err != nil {
		fmt.Println("Error on initializing listener:", err)
		os.Exit(1)
	}

	defer listener.Close()

	fmt.Println("Listening on - ", *addr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	fmt.Println(
		time.Now().UTC().Format("2006-01-02 15:04:05"),
		": Hanlding connection from - ",
		conn.RemoteAddr())

	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	reqLen, err := conn.Read(buf)
	var ipAddress string
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	} else {
		ipAddress = strings.TrimSpace(string(buf[:reqLen]))
		fmt.Println("get request:", ipAddress)
	}

	reply := checkIPExist(ipAddress)

	conn.Write([]byte(reply))
	fmt.Println("send reply:", reply)
	conn.Close()
}

func checkIPExist(ip string) string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			ipString := ipnet.IP.To4().String()
			if ip == ipString {
				return "1"
			}
		}
	}

	return "0"
}
