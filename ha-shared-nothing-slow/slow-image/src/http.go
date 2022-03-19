package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func initHttpServer(server *http.Server) {
	defer notifyTermination()

	http.HandleFunc("/set/", httpSetHandler)
	http.HandleFunc("/", httpDefaultHandler)

	log.Printf("http server listening on %s...", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Print("web server graceful shutdown")
			return
		}
		log.Fatal(err)
	}
}

func httpSetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	path := r.URL.Path
	if len(path) < len("/set/") {
		http.Error(w, "invalid URI", http.StatusInternalServerError)
		return
	}
	parts := strings.Split(path[len("/set/"):], "/")
	if len(parts) < 2 {
		http.Error(w, "invalid URI", http.StatusInternalServerError)
		return
	}

	key := strings.ToLower(parts[0])
	val, err := strconv.Atoi(parts[1])
	if err != nil {
		http.Error(w, fmt.Sprintf("%s could not be converted to an integer: %v", parts[1], err), http.StatusInternalServerError)
		return
	}
	switch key {
	case "streamdelay":
		if err := config.setStreamDelay(val); err != nil {
			http.Error(w, fmt.Sprintf("error setting stream delay: %v", err), http.StatusInternalServerError)
			return
		}
	case "conndelay":
		if err := config.setConnDelay(val); err != nil {
			http.Error(w, fmt.Sprintf("error setting connection delay: %v", err), http.StatusInternalServerError)
			return
		}
	case "buffersize":
		if err := config.setBufferSize(val); err != nil {
			http.Error(w, fmt.Sprintf("error setting write buffer size: %v", err), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "unrecognized key", http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "ok")
}

func httpDefaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "number of connections=%d\n", registry.size())
	fmt.Fprintln(w, &config)
}
