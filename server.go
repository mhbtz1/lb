package main

import (
	"os"
	"net/http"
	"fmt"
	//"html/template"
	"strings"
)

type Page struct {
	Title string
	Body []byte
}

func (p *Page) save() error {
    filename := p.Title + ".txt"
    return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
    filename := title + ".txt"
    body, _ := os.ReadFile(filename)
    return &Page{Title: title, Body: body}, nil
}

func extractHostAndPort(r *http.Request) (string, string, error) {
	// Split the Host field into host and port
	parts := strings.Split(r.Host, ":")
	if len(parts) == 1 {
		// No port specified, use the default HTTP or HTTPS port
		return parts[0], "", nil
	} else if len(parts) == 2 {
		// Host and port specified
		return parts[0], parts[1], nil
	} else {
		// Unexpected format
		return "", "", fmt.Errorf("unexpected Host format: %s", r.Host)
	}
}
func defaultHandler(w http.ResponseWriter, r *http.Request){
	_, port, _ := extractHostAndPort(r)
	fmt.Printf("Hello from web server %s.\n", port)
	/*
	t, _ := template.ParseFiles(fmt.Sprintf("./%s.html"), port)
  t.Execute(w, r)	
	*/
}

