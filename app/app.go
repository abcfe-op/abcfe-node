package app

import (
	"flag"
	"fmt"

	conf "github.com/abcfe/abcfe-node/config"
)

var port = flag.Int("port", 3000, "Set port of the server")
var mode = flag.String("mode", "rest", "Choose between 'boot' and 'validator'")
var configPath = flag.String("config", "", "Set path of the config file")

type App struct {
	stop chan struct{}
	Conf conf.Config
}

func New() (*App, error) {
	flag.Parse()
	cfg, err := conf.NewConfig(*configPath)
	if err != nil {
		fmt.Println("Failed to initialized application: ", err)
		return nil, err
	}

	app := &App{
		stop: make(chan struct{}),
		Conf: *cfg,
	}

	return app, nil
}

func (a *App) Wait() {
	<-a.stop // 채널에서 값 읽으려고 시도
}

func (a *App) Terminate() {
	close(a.stop)
}
