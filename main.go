package main

import (
	"fmt"
	"proxyServer/proxy"

	"net/http"
)

func main() {
	var p proxy.Proxy
	p.NewProxy("memes")



	fmt.Println("Starting server on localhost:8080")
	addr := ":8080"
	err := http.ListenAndServe(addr, &p)
	if err != nil {
		fmt.Println("error of starting server")
	}

}
