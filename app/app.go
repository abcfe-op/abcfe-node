package app

import (
	"flag"
	"fmt"

	"github.com/abcfe-op/abcfe-node/cli"
	"github.com/abcfe-op/abcfe-node/config"
	"github.com/abcfe-op/abcfe-node/db"
	"github.com/abcfe-op/abcfe-node/rpc"

	log "github.com/abcfe-op/abcfe-node/common/logger"
)

var port = flag.Int("port", 4000, "Set port of the server")
var mode = flag.String("mode", "rest", "Choose between 'auto' and 'rest'")
var configPath = flag.String("config", "../config/config.toml", "Set path of the config file")

type App struct {
	stop chan struct{}
	cfg  *config.Config
}

func New() (*App, error) {
	flag.Parse()
	cfg := config.NewConfig(*configPath)

	r := &App{
		stop: make(chan struct{}, 1),
		cfg:  cfg,
	}

	log.InitLogger(r.cfg)
	log.Info(fmt.Sprintf("Starting %s at Port: %d Mode: %s", r.cfg.Common.ServiceName, *port, *mode))

	defer db.Close()
	db.InitDB()

	go rpc.Start(*port)
	cli.Start(*port, *mode)

	return &App{}, nil
}

func (p *App) Wait() {
	<-p.stop
}

func (p *App) Terminate() {
	defer close(p.stop)
}
