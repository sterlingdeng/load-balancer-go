package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// LoadBalancer is the configuration file for the loadbalancer
type LoadBalancer struct {
	Host      string
	Port      string
	Scheme    string
	Servers   []Server
	getServer func() *Server
}

func (lb *LoadBalancer) getURL() string {
	return lb.Scheme + "://" + lb.Host + ":" + lb.Port
}

func (lb *LoadBalancer) initialize() {
	lb.setDefaultParams()
	lb.getServer = lb.initIterator()
}

func (lb *LoadBalancer) initIterator() func() *Server {
	currServerIdx := 0
	serverCount := len(lb.Servers) - 1
	return func() *Server {
		server := &lb.Servers[currServerIdx]
		currServerIdx++
		if currServerIdx > serverCount {
			currServerIdx = 0
		}
		return server
	}
}

func (lb *LoadBalancer) handleRequest(w http.ResponseWriter, r *http.Request) {
	lb.tryServers(w, r)
}

func (lb *LoadBalancer) tryServers(w http.ResponseWriter, r *http.Request) {
	// select the server
	server := lb.getServer()
	err := server.checkHealth()

	count := 1

	for err != nil && count < len(lb.Servers) {
		server = lb.getServer()
		err = server.checkHealth()
		count++
	}

	if err != nil {
		log.Fatal("[ERROR] All downstream connections failed")
	}

	server.proxyRequest(w, r, lb)
}

func (s *Server) proxyRequest(w http.ResponseWriter, r *http.Request, lb *LoadBalancer) {
	uri, err := url.Parse(s.getURL() + r.RequestURI)

	if err != nil {
		fmt.Println("proxyRequest error: ", err)
	}

	// when forwarding, need to set X-Forwarded-Host
	r.URL = uri
	r.Header.Set("X-Forwarded-Host", r.Host)
	r.Header.Set("Origin", lb.getURL())
	// r.Host is the address of the server
	r.Host = s.getURL()
	r.RequestURI = ""

	client := &http.Client{}
	resp, err := client.Do(r)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	for k, v := range resp.Header {
		fmt.Println("header kv", k, " ", v)
		w.Header().Set(k, strings.Join(v, ";"))
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func (s *Server) checkHealth() error {
	resp, err := http.Get(s.getURL())
	if err != nil {
		fmt.Println("health check failed: ", err)
		return err
	}
	if resp.StatusCode >= 400 {
		fmt.Println("health check failed: ", err)
		return err
	}
	return nil
}
