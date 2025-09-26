I think it's still important to think even though we have the generative AI now. Thinking always has been and always will be good and useful. Before we can start thinking about something more deeply we have to learn and understand the basics. And one of the best ways how to learn anything is by doing it. Playing is a way of doing. You should probably not play in production though.

Take TLS for example. It used to be called SSL before and it's the protocol that secures the network communication. It sits between Application and Transport layers:

```
TCP/IP Layer    | Protocol / Medium
--------------- | ----------------------
Application     | HTTP, SMTP, DNS, ...
Transport       | TCP, UDP
Internet        | IP
Link / Network  | Ethernet, WiFi
Physical        | Cables, radio, optical
```

(This is the TCP/IP networking model which is simpler than the OSI model. However, if you like pizza and sausages, the OSI model is easy to remember: Please Do Not Throw The Sausage Pizza Away :-)

## TCP

TLS wraps application data (like HTTP) before it goes over TCP. Therefore let's have a look at how TCP works first. We build a simple TCP server:

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

The data (`hello\n`) goes over the network in plaintext. If someone eavesdrops on the network, for example using Wireshark, they see the transferred data:

<img width="594" height="65" alt="image" src="https://github.com/user-attachments/assets/b3e1daa4-fd69-49a8-ac85-399f10e50f51" />

## TLS

TLS is encrypting the application data so that they are not readable when someone gets them from the network. Now we use the [crypto/tls](https://pkg.go.dev/crypto/tls) standard library package instead of [net](https://pkg.go.dev/net). Also we need to supply the TLS certificate and private key:

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

I added logging to the `echo` function so we can see what's going on with the connection:

```go
// ./tls/server.go
func echo(conn net.Conn) {
	n, err := io.Copy(conn, conn)
	log.Printf("sent %d bytes to %s, err: %v", n, conn.RemoteAddr(), err)
	conn.Close()
}
```

I used the [mkcert](https://github.com/FiloSottile/mkcert) tool to create certificate and key file for localhost:

```
$ mkcert localhost
```

But when we send some data to the server now, we don't see anything echoed back:

```
$ go run ./tls/server.go &
$ echo hello | nc localhost 4321
```

Let's have a look at the server logs:

```
$ go run main.go 
2025/09/25 18:00:06 sent 0 bytes to 127.0.0.1:52130, err: tls: first record does not look like a TLS handshake
```

Yes, netcat (`nc`) can't speak TLS: the yellow steps. It can only speak TCP: the blue steps (the standard three-way TCP handshake): 

<img width="542" height="351" alt="image" src="https://github.com/user-attachments/assets/b567ea3b-0c35-4f40-bcb0-491af382f403" />

But `openssl` can:

```
$ echo hello | openssl s_client -connect localhost:4321 -servername localhost -quiet 2> /dev/null 
hello
```

If we capture the data from the network now, we can't read it since it's encrypted:

<img width="649" height="219" alt="image" src="https://github.com/user-attachments/assets/eb4da1a3-c077-4753-be32-8892fd595509" />

## HTTP

TLS is often used to secure the HTTP protocol. HTTP on top of TLS becomes HTTPS.

So let's start with HTTP. Here's the HTTP version of the echo server from the above:

```go
// ./http/server.go
func main() {
	http.HandleFunc("/", echo)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func echo(resp http.ResponseWriter, req *http.Request) {
	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(resp,
			"Error reading request body",
			http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	resp.Write(reqBody)
}
```

Let's start the server and send some data to it via `curl`:

```
$ go run http/server.go &
$ curl localhost:8080 --data hello
hello
```

We see it gets echoed back. In plaintext, unencrypted.

## HTTPS

To secure our HTTP communication we just need to use `http.ListenAndServeTLS` instead of `http.ListenAndServe` (both are part the powerful [net/http](https://pkg.go.dev/net/http) stdlib package). And we need to supply files containing a certificate and matching private key:

```go
// ./https/server.go
http.HandleFunc("/", echo)
log.Fatal(http.ListenAndServeTLS("localhost:4430",
        "localhost.pem", "localhost-key.pem", nil))
```

On the client side, we need to specify `https://` in the URL (`curl` defaults to `http://`) and port `4430` instead of `8080`. Also we either skip the server certificate verification by using `--insecure` (or `-k`) or we supply the server certificate file via `--cacert`.

```
$ go run https/server.go &
$ curl https://localhost:4430 --data hello --insecure
hello
$ curl https://localhost:4430 --data hello --cacert localhost.pem 
hello
```