package main

import (
	"encoding/json"
	"net/http"
)

func getLastFiveMessages(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		return
	}

	historyMux.Lock()
	defer historyMux.Unlock()

	count := len(messageHistory)
	if count > 5 {
		count = 5
	}

	lastFive := messageHistory[len(messageHistory)-count:]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lastFive)
}
