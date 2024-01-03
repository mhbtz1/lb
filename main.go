package main

import (
	"fmt"
	"net/http"
	"sync"
	"load-balancer/concqueue"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Server struct {
		Port int `yaml:"Port"`
		Address string `yaml:"Address"`
}

type ServerConfig struct {
	Server1 Server `yaml:"Server1"`
	Server2 Server `yaml:"Server2"`
	
	LoadBalancer struct {
		Port int `yaml:"Port"`
		Address string `yaml:"Address"`
	} `yaml:"LoadBalancer"`

}


var conc *concqueue.ConcurrentQueue
var config ServerConfig

func main() {
	// Start the first HTTP server on port 5000
	var wg sync.WaitGroup
	
	conc = concqueue.MakeQueue(100)
	
	data, _ := ioutil.ReadFile("server_ports.yml")
	err := yaml.Unmarshal(data, &config)
	
	if (err != nil) {
		return
	}
		
	wg.Add(1)
	go startServer(config.Server1.Port, defaultHandler, &wg)

	// Start the second HTTP server on port 8008
	wg.Add(1)
	go startServer(config.Server2.Port, defaultHandler, &wg)

	// Start the third HTTP server on port 8080
	wg.Add(1)
	go startServer(config.LoadBalancer.Port, loadHandler, &wg)

	// Wait indefinitely to keep the program running
	wg.Wait()
}

func startServer(port int, h func(w http.ResponseWriter, r *http.Request), wg *sync.WaitGroup) {
	defer wg.Done()
	
	mux := http.NewServeMux()
	// Register a simple handler for the server
	mux.HandleFunc("/", h)

	// Start the server on the specified port
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Server listening on %s...\n", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		fmt.Printf("Error starting server on port %d: %v\n", port, err)
	}
}
