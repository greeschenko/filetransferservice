package main

import (
    "fmt"
    "net/http"
)

func main() {
    fmt.Println("DEMO SERVER SRARTED!!!")
    http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "Hello World!")
    })
    http.ListenAndServe(":8080", nil)
}

