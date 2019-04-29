package main

import (
	"fmt"
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

const ()

// Server contains information on the server of which to balance load over
type Server struct {
	Name   string
	Scheme string
	Host   string
	Port   string
}

func (s *Server) getURL() string {
	return s.Scheme + "://" + s.Host + ":" + s.Port
}

func (lb *LoadBalancer) setDefaultParams() {
	if lb.Host == "" {
		lb.Host = "localhost"
	}
	if lb.Port == "" {
		lb.Port = "3000"
	}
	if lb.Scheme == "" {
		lb.Scheme = "http"
	}
}

// ParseConfig parses config.yml
func ParseConfig() *LoadBalancer {

	lb := &LoadBalancer{}

	file, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		fmt.Println("Error reading config.yml")
	}

	err = yaml.Unmarshal(file, lb)
	if err != nil {
		fmt.Println("Error parsing config.yml")
	}

	return lb
}
