package server

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func myHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("------")

	fmt.Printf("Host:             %s\n", r.Host)
	fmt.Printf("Header:           %s\n", r.Header)
	fmt.Printf("RequestURI:       %s\n", r.RequestURI)
	// fmt.Printf("Body:             %s\n", r.Body)
	fmt.Printf("Form:             %s\n", r.Form)
	fmt.Printf("Method:           %s\n", r.Method)
	fmt.Printf("RemoteAddr:       %s\n", r.RemoteAddr)
	fmt.Printf("TransferEncoding: %s\n", r.TransferEncoding)

	APIRoute(w, r)
}

// Run ...
func Run(port int, args []string) error {

	s := &http.Server{
		Addr:           ":" + strconv.Itoa(port),
		Handler:        http.HandlerFunc(myHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err := s.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}
