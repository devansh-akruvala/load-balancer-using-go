package main

import (
	"fmt"
	"net/http"
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

type Server interface {
	Address() string
	IsAlive() bool
	serve(resW http.ResponseWriter, reqW *http.Request)
}

type LoadBalancer struct {
	port            string
	roundRobinCount int
	servers         []Server
}

func newLoadBalancer(port string, servers []Server) *LoadBalancer {
	return &LoadBalancer{
		port:            port,
		roundRobinCount: 0,
		servers:         servers,
	}
}

func main() {

}
