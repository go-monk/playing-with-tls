package main

import (
	"io"
	"log"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", "localhost:1234")
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
	io.Copy(conn, conn)
	conn.Close()
}
