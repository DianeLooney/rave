package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

const SockAddr = "/tmp/rave.sock"

func catServer(conn net.Conn) {
	defer conn.Close()

	var a, b, c string
	n, err := fmt.Fscanln(conn, &a, &b, &c)
	fmt.Printf("fmt.FScanln: %v, %v\n", n, err)
	fmt.Printf("%s\n%s\n%s\n", a, b, c)
}

func main() {
	if err := os.RemoveAll(SockAddr); err != nil {
		log.Fatal(err)
	}

	l, err := net.Listen("unix", SockAddr)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	defer l.Close()

	for {
		// Accept new connections, dispatching them to echoServer
		// in a goroutine.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
		}

		go catServer(conn)
	}
}
