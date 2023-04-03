package main

import (
	"fmt"
	data "github.com/anthonynixon/link-shortener-backend/internal/cloud"
	"github.com/anthonynixon/link-shortener-backend/internal/handlers/link"
	"github.com/anthonynixon/link-shortener-backend/internal/router"
	"log"
	"os"
)

var PORT = ""

func init() {
	PORT = os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	data.Initialize()
}

func main() {
	router := router.New()

	// Add Routes
	link.AddLinkV1(router)

	log.Printf("Running link-shortener-backend on :%s...", PORT)
	err := router.Run(fmt.Sprintf(":%s", PORT))
	if err != nil {
		log.Fatal(err.Error())
	}
}
