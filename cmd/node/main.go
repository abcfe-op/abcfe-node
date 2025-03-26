package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/abcfe/abcfe-node/app"
	"github.com/abcfe/abcfe-node/common/logger"
	"github.com/abcfe/abcfe-node/storage"
)

func main() {
	app, err := app.New()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
		os.Exit(1)
	}

	if err := logger.InitLogger(&app.Conf); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
		os.Exit(1)
	}

	db, err := storage.InitDB(&app.Conf)
	if err != nil {
		logger.Error("Failed to load db: ", err)
	}

	fmt.Println("temp db log: ", db)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM) // OS 시그널을 채널로 전달
	go func() {
		sig := <-sigCh
		logger.Info("Arrived terminate signal: ", sig)
		app.Terminate()
	}()

	logger.Info("Node start.")
	app.Wait()
	logger.Info("Node terminated.")
}
