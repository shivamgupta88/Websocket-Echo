package main

import "net/http"

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", handleWebSocket)
	mux.HandleFunc("/history", getLastFiveMessages)
	return mux
}
