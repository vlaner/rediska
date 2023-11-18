package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/vlaner/rediska/internal/server"
)

func main() {
	s := server.New(":3000")
	ctx, cancel := context.WithCancel(context.Background())
	exit := make(chan os.Signal, 1)
	defer cancel()

	if err := s.Start(ctx); err != nil {
		log.Fatalln(err)
	}

	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	<-exit

	log.Println("stopping server...")
	s.Stop()
	log.Println("server stopped")
}
