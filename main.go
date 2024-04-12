package main

import (
	"fmt"
	"net/http/httputil"
	"net/url"
	"os"
)

type server struct {
	address string
	proxy   *httputil.ReverseProxy
}

func newServer(addr string) *server {
	serverUrl, err := url.Parse(addr)
	handleErr(err)

	return &server{
		address: addr,
		proxy:   httputil.NewSingleHostReverseProxy(serverUrl),
	}
}

func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {

}
