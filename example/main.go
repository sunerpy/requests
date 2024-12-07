package main

import (
	"fmt"
	"log"

	"github.com/sunerpy/requests"
	"github.com/sunerpy/requests/url"
)

func main() {
	// Enable HTTP/2 globally
	requests.SetHTTP2Enabled(true)
	log.Println("HTTP/2 enabled globally")
	// Perform a GET request
	params := url.NewURLParams()
	params.Set("query", "golang")
	resp, err := requests.Get("https://httpbin.org/get", params)
	if err != nil {
		log.Fatalf("GET Error: %v", err)
	}
	log.Printf("GET Response Status Code: %d\n", resp.StatusCode)
	log.Printf("GET Response Protocol: %s\n", resp.Proto)
	fmt.Println("GET Response Text:", resp.Text())
	// Perform a GET request using a session with HTTP/1.1
	session := requests.NewSession().WithHTTP2(false)
	defer session.Close()
	log.Println("Using a new Session with HTTP/1.1")
	req, err := requests.NewRequest("GET", "https://httpbin.org/get", params, nil)
	if err != nil {
		log.Fatalf("Request creation error: %v", err)
	}
	respHTTP1, err := session.Do(req)
	if err != nil {
		log.Fatalf("HTTP/1.1 GET Error: %v", err)
	}
	log.Printf("HTTP/1.1 GET Response Protocol: %s\n", respHTTP1.Proto)
	fmt.Println("HTTP/1.1 GET Response Text:", respHTTP1.Text())
	// Perform a POST request
	form := url.NewForm()
	form.Set("username", "john_doe")
	respPost, err := requests.Post("https://httpbin.org/post", form)
	if err != nil {
		log.Fatalf("POST Error: %v", err)
	}
	log.Printf("POST Response Protocol: %s\n", respPost.Proto)
	fmt.Println("POST Response Text:", respPost.Text())
	fmt.Println("POST Response Headers:", respPost.Headers)
	fmt.Println("POST Response Cookies:", respPost.Cookies)
	fmt.Println("POST Response Status Code:", respPost.StatusCode)
}
