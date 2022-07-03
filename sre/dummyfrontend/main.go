package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"
	"fmt"
)

func respondRandomly(w http.ResponseWriter, r *http.Request) {	
	fmt.Fprintf(w, "Hello, %s", r.URL.Path[1:])
}

func main() {
	rand.Seed(time.Now().UnixNano())
	http.HandleFunc("/", respondRandomly)
	err := http.ListenAndServe(":3001", nil)
	if err != nil {
		log.Fatal(err)
	}
}
