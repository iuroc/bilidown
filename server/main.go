package main

import (
	"bilidown/router"
	"fmt"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("../client/dist")))
	http.Handle("/api/", http.StripPrefix("/api", router.API()))

	fmt.Println("http://127.0.0.1:8098")
	http.ListenAndServe("127.0.0.1:8098", nil)
}
