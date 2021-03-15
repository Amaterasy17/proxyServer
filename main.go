package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"proxyServer/proxy"
)

func main() {
	var p proxy.Proxy
	p.NewProxy("config")

	var pem, key, proto string
	flag.StringVar(&pem, "pem", "RootCA.pem", "")
	flag.StringVar(&key, "key", "RootCA.key", "")
	flag.StringVar(&proto, "proto", "http", "")
	flag.Parse()

	if proto != "http" && proto != "https" {
		log.Println("Protocol must be either http or https")
		return
	}


	fmt.Println("Starting server on localhost:8080")
	addr := ":8080"
	server := &http.Server{
		Addr: addr,
		Handler: &p,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	//
	if proto == "http" {
		log.Fatal(server.ListenAndServe())
	} else {
		log.Fatal(server.ListenAndServeTLS(pem, key))
	}

}

