package main

import (
	"net/http"
	"fmt"
	"math/rand"
	"net/url"
	"net/http/httputil"
	"time"
	"load-balancer/concqueue"
	"sync"
)

var mtx sync.Mutex

func handleRequest(destinationURL string, w http.ResponseWriter, r *http.Request) (*http.Request, *url.URL) {
	destinationURI, _ := url.Parse(destinationURL)
	fmt.Printf("Destination URI: %s\n", destinationURI)
	
	newRequest, _ := http.NewRequest(r.Method, destinationURI.String(), r.Body)
	for key, values := range r.Header {
		for _, value := range values {
			newRequest.Header.Add(key, value)
		}
	}


	return newRequest, destinationURI
}


func isServerUp(endpoint string) bool {
	client := http.Client{
		Timeout: 2 * time.Second, // Adjust the timeout as needed
	}
	
	req, err := http.NewRequest("GET", endpoint, nil)
	if (err != nil ){
		fmt.Println("Error forming request.")
		return false
	}

	req.Header.Set("ignore-request", "true")

	resp, err := client.Do(req)

	if err != nil {
		return false
	}
	// Check the response status code
	if resp.StatusCode == http.StatusOK {
		return true
	}

	return false
}


/*
Use concurrent queue for potentially paarallel requests made to one load balancer 
that cannot be serviced immediately
*/
func loadHandler( w http.ResponseWriter, r *http.Request ){
	
	mtx.Lock()
	val := r.Header.Get("ignore-request")
	if val == "" {
		conc.Enqueue(r)
	} 
	
	fn := func ( w http.ResponseWriter, r *http.Request, q *concqueue.ConcurrentQueue) {
		var req *http.Request

		
		server_ports := []int{config.Server1.Port, config.Server2.Port}
		fmt.Printf("Current queue length: %d\n", q.CheckSize() )
		
		if (q.CheckSize() == 0){
			return;
		}
		req = q.Dequeue()
		
		port := balance(server_ports)
		newRequest, destinationURI := handleRequest(fmt.Sprintf("http://127.0.0.1:%d", port), w, req);
		
		if ( newRequest != nil ){
			reverseProxy := httputil.NewSingleHostReverseProxy(destinationURI) 
			reverseProxy.ServeHTTP(w, newRequest)
		}
		fmt.Println("------------------------------------------------------------------")
		
	}

	fn(w, r, conc)
	mtx.Unlock()
}


func balance(server_ports []int) int {
	var active_servers []int 

	for _, v := range server_ports {
			endpoint := fmt.Sprintf("http://localhost:%d", v) 
			if ( isServerUp(endpoint) ){
				fmt.Printf("Server %d is up!\n", v)
				active_servers = append(active_servers, v)
			} 
	}

	 rand.Seed(time.Now().UnixNano())
	 index := rand.Intn(len(active_servers))
	 return active_servers[index]
}

func round_robin(server_ports []int) int {
	return 0
}
