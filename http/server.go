package main

import (
	"io"
	"log"
	"net/http"
)

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
