package main

import "log"

func main() {
	a := NewApp()
	err := a.Start()
	if err != nil {
		log.Fatalf("Error while starting app: %v", err)
	}
	sig := <-a.shutdown
	log.Printf("shutdown signal received: %v\n", sig)
	err = a.Stop()
	if err != nil {
		log.Fatalf("Error while shutting down app: %v", err)
	}
}
