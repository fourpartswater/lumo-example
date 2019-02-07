package main

import (
	"fmt"
	"net/http"

	"github.com/unchartedsoftware/lumo-example/api/mandelbrot"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/mandelbrot/", mandelbrot.Handler)
	mux.Handle("/", http.FileServer(http.Dir("./build")))
	fmt.Println("Server listening on port 8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println(err)
	}
}
