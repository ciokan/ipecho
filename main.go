package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type response struct {
	Error        string `json:"error"`
	Ip           string `json:"ip"`
	Port         string `json:"port"`
	ForwardedFor string `json:"forwarded_for"`
}

func getIp(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	ip, port, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		if jerr := json.NewEncoder(w).Encode(response{Error: err.Error()}); jerr != nil {
			log.Fatalf("failed to encode json")
		}
	}

	forward := req.Header.Get("X-Forwarded-For")

	display := req.URL.Query().Get("display")

	if display == "simple" {
		switch req.URL.Query().Get("field") {
		case "port":
			fmt.Fprintf(w, port)
			break
		case "forwarded-for":
			fmt.Fprintf(w, forward)
			break
		default:
			fmt.Fprintf(w, ip)
		}

		return
	}

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

	r.GET("/", getIp)

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.Serve(l, r))
}
