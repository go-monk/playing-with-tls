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
func main() {
        ln, err := net.Listen("tcp", "localhost:8080")
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
$ echo hello | nc localhost 8080
hello
```

The data (`hello\n`) goes over the network in plaintext. If someone eavesdrops the connection, for example using Wireshark, he sees the data:

<img width="594" height="65" alt="image" src="https://github.com/user-attachments/assets/b3e1daa4-fd69-49a8-ac85-399f10e50f51" />

## TLS

## HTTP

TLS is most often used to secure the HTTP protocol. HTTP on top of TLS becomes HTTPS.

## HTTPS
