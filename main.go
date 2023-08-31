package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load(".env")

	fmt.Println("Init")
	portString := os.Getenv("PORT")

	if portString == "" {
		log.Fatal("Port not found")
	}

	fmt.Println("Port: $s", portString)
}