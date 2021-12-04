package main

import (
	"net/http"

	"github.com/Chiorufarewerin/gitchat/function"
	"github.com/Chiorufarewerin/gitchat/internal/environment"
)

func main() {
	http.HandleFunc("/", function.Handler)

	http.ListenAndServe(environment.Addr, nil)
}
