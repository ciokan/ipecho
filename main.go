package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"net"
	"net/http"
)

type response struct {
	Error        string `json:"error"`
	Ip           string `json:"ip"`
	Port         string `json:"port"`
	ForwardedFor string `json:"forwarded_for"`
}

func getIP(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	
	ip, port, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		if jerr := json.NewEncoder(w).Encode(response{Error: err.Error()}); jerr != nil {
			panic(jerr)
		}
	}
	
	forward := req.Header.Get("X-Forwarded-For")
	
	if jerr := json.NewEncoder(w).Encode(response{
		Ip:           ip,
		Port:         port,
		ForwardedFor: forward,
	}); jerr != nil {
		panic(jerr)
	}
}

func main() {
	r := httprouter.New()
	
	r.GET("/", getIP)
	
	l, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatal(err)
	}
	
	log.Fatal(http.Serve(l, r))
}
