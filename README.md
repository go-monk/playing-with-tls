I think it's still important to think even though we have the generative AI now. Thinking always have been and always will be good and useful. Think about it :-). Before you can think about something you have to learn and understand the basics. And one of the best ways how to learn something is by doing it.

Take for example TLS. It used to be called SSL before and it's the protocol that secures the Internet. To be more precise it usually secures the HTTP traffic. It sits between Application and Transport layers:

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
