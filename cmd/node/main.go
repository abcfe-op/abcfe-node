package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/abcfe/abcfe-node/app"
	"github.com/abcfe/abcfe-node/common/logger"
)

func main() {
	nodeApp, err := app.New()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
		os.Exit(1)
	}

	if err := logger.InitLogger(&nodeApp.Conf); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
		os.Exit(1)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM) // OS 시그널을 채널로 전달
	go func() {
		sig := <-sigCh
		logger.Info("Arrived terminate signal: ", sig)
		nodeApp.Terminate()
	}()

	logger.Info("Node start.")
	nodeApp.Wait()
	logger.Info("Node terminated.")
}
