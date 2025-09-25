I think it's still important to think even though we have the generative AI now. Thinking always have been and always will be good and useful. Think about it :-). Before you can think about something you have to learn and understand the basics. And one of the best ways how to learn something is by doing it or in other words playing with it.

Take for example TLS. It used to be called SSL before and it's the protocol that secures the network communication. It sits between Application and Transport layers:

```
TCP/IP Layer        | Example / Role
------------------- | -----------------------
Application         | HTTP, SMTP, DNS, SSH
Transport           | TCP / UDP
Internet            | IP
Link / Network      | Ethernet, WiFi
Physical            | Cables, radio, optical
```

(This the TCP/IP networking model which is simpler the the OSI model. However, if you like pizza the OSI model is easy to remember: Please Do Not Throw The Sausage Pizza Away :-)

## TCP

TLS wraps application data before it goes over TCP. Let's have a look at how TCP works first. Let's build a simple TCP server:

```go
// ./tcp/server.go
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
```

The server echoes back anything a client sends to it:

```
$ go run ./tcp/server.go &
$ echo hello | nc localhost 1234
hello
```

The data (`hello\n`) goes over the network in plaintext. If someone eavesdrops the connection, for example using Wireshark, he sees the data:

<img width="594" height="65" alt="image" src="https://github.com/user-attachments/assets/b3e1daa4-fd69-49a8-ac85-399f10e50f51" />

## TLS

TLS is encrypting the application data so that they are not readable when someone reads them from the network. Now we use the `crypto/tls` standard library package instead of `net`. Also we need to supply the TLS certificate and private key:

```go
// ./tls/server.go
cert, err := tls.LoadX509KeyPair("localhost.pem", "localhost-key.pem")
if err != nil {
        log.Fatal(err)
}

config := &tls.Config{
        Certificates: []tls.Certificate{cert},
}

ln, err := tls.Listen("tcp", "localhost:4321", config)
// The rest of the code is as above...
```

Also I added logging so we can see what's going on with the connection:

```go
func echo(conn net.Conn) {
	n, err := io.Copy(conn, conn)
	log.Printf("sent %d bytes to %s, err: %v", n, conn.RemoteAddr(), err)
	conn.Close()
}
```

I created the certificate a key file for localhost using the [mkcert](https://github.com/FiloSottile/mkcert) tool:

```
$ mkcert localhost
```

When we send some data to the server now, we don't see anything echoed back:

```
$ go run ./tls/server.go &
$ echo hello | nc localhost 4321
```

If we look at the server logs we start to see why:

```
$ go run main.go 
2025/09/25 18:00:06 sent 0 bytes to 127.0.0.1:52130, err: tls: first record does not look like a TLS handshake
```

Netcat (`nc`) can't speak TLS - the yellow steps; it can only speak TCP - the blue steps (the standard three-way TCP handshake): 

<img width="542" height="351" alt="image" src="https://github.com/user-attachments/assets/b567ea3b-0c35-4f40-bcb0-491af382f403" />

But `openssl` can:

```
$ echo hello | openssl s_client -connect localhost:4321 -servername localhost -quiet 2> /dev/null 
hello
```

## HTTP

TLS is most often used to secure the HTTP protocol. HTTP on top of TLS becomes HTTPS.

## HTTPS
