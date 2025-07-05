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

	// rest api start
	if err := app.NewRest(); err != nil {
		logger.Error("Failed to start services:", err)
		app.Terminate()
		os.Exit(1)
	}
	// json-rpc start

	// grpc start

	// PoS start
	// config 기준으로 역할 규정 (root, validator, sentry)

	app.Wait()
	logger.Info("Node terminated.")
}
