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

func (s *server) Address() string {
	return s.address
}

func (s *server) IsAlive() bool {
	return true
}

func (s *server) serve(resW http.ResponseWriter, reqW *http.Request) {
	s.proxy.ServeHTTP(resW, reqW)
}

func (lb *LoadBalancer) IsAlive(s *server) bool {
	return true
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

func (lb *LoadBalancer) getNextAvailableServer() Server {
	fmt.Println(lb.roundRobinCount)
	servers := lb.servers[lb.roundRobinCount%len(lb.servers)]
	for !servers.IsAlive() {
		lb.roundRobinCount++
		servers = lb.servers[lb.roundRobinCount%len(lb.servers)]
	}
	lb.roundRobinCount++
	fmt.Printf(" after %d", lb.roundRobinCount)
	return servers

}

func (lb *LoadBalancer) serve(resW http.ResponseWriter, req *http.Request) {
	targetServer := lb.getNextAvailableServer()
	fmt.Printf("forwarding request to address %q\n", targetServer.Address())
	targetServer.serve(resW, req)
}

func main() {
	servers := []Server{
		newServer("https://www.google.com/"),
		newServer("https://www.bing.com/"),
		newServer("https://duckduckgo.com/"),
	}

	lb := newLoadBalancer("8000", servers)
	handleRedirect := func(resW http.ResponseWriter, req *http.Request) {
		lb.serve(resW, req)
	}
	http.HandleFunc("/", handleRedirect)

	fmt.Printf("Server is Listening at 'localhost:%s'\n", lb.port)
	http.ListenAndServe(":"+lb.port, nil)
}
