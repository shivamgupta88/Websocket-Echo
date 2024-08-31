package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	routes := SetupRoutes()
	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", routes))
}
