package main

import (
	"flag"
	"log"
	"os"

	"github.com/abcfe/abcfe-node/app"
	"github.com/abcfe/abcfe-node/common/logger"
)

func main() {
	flag.Parse()

	app, err := app.New()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
		os.Exit(1)
	}

	app.SigHandler()

	logger.Info("Node start.")
	app.Wait()
	logger.Info("Node terminated.")
}
