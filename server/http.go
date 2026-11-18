package server

import (
	"io"
	"log"
	"net/http"
)

func Health() {

	healthHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "ok\n")
	}

	http.HandleFunc("/health", healthHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
