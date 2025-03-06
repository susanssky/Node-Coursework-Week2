package main

import (
	"log"
	"node-week-02-with-chi/api"
)

// @title CYF Chat Application API
// @version 1.0
// @description This is a RESTful API for the CYF chat application, providing message management capabilities.
// @host localhost:4001
// @BasePath /api/v1
func main() {
	server := api.NewAPIServer(":4001")

	err := server.Run()
	if err != nil {
		log.Fatalf("Server error:%v", err)
	}

}
