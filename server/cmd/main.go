package main

import (
	"log"
)

func main() {
	a := NewApp()
	a.Start()
	sig := <-a.shutdown
	log.Printf("shutdown signal received: %v\n", sig)
	if err := a.Stop(); err != nil {
		log.Fatalf("Failed to gracefully shutdown: %v", err)
	}
	log.Println("Shutdown complete!")
}
