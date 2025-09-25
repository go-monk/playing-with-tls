package main

import (
	"crypto/tls"
	"io"
	"log"
	"net"
)

func main() {
	cert, err := tls.LoadX509KeyPair("localhost.pem", "localhost-key.pem")
	if err != nil {
		log.Fatal(err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	ln, err := tls.Listen("tcp", "localhost:4321", config)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go echo(conn)
	}
}

func echo(conn net.Conn) {
	n, err := io.Copy(conn, conn)
	log.Printf("sent %d bytes to %s, err: %v", n, conn.RemoteAddr(), err)
	conn.Close()
}
