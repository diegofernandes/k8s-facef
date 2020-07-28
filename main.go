package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

var port = flag.Int("port", 8090, "listen port")

func greet(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("Bye World! %s", time.Now())

	fmt.Fprint(w, msg)
}

func main() {
	flag.Parse()
	addr := fmt.Sprintf(":%v", *port)
	log.Printf("Starting hello api on %s", addr)
	http.HandleFunc("/", greet)
	log.Println(http.ListenAndServe(addr, nil))
}
