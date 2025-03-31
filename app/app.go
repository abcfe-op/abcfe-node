package app

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/abcfe/abcfe-node/common/logger"
	conf "github.com/abcfe/abcfe-node/config"
	"github.com/abcfe/abcfe-node/storage"
	"github.com/syndtr/goleveldb/leveldb"
)

var port = flag.Int("port", 3000, "Set port of the server")
var mode = flag.String("mode", "rest", "Choose between 'boot' and 'validator'")
var configPath = flag.String("config", "", "Set path of the config file")

type App struct {
	stop chan struct{}
	Conf conf.Config
	DB   *leveldb.DB // db내의 mutex는 복사되면 안됨
}

func New() (*App, error) {
	cfg, err := conf.NewConfig(*configPath)
	if err != nil {
		fmt.Println("Failed to initialized application: ", err)
		return nil, err
	}

	if err := logger.InitLogger(cfg); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
		os.Exit(1)
	}

	db, err := storage.InitDB(cfg)
	if err != nil {
		logger.Error("Failed to load db: ", err)
	}

	app := &App{
		stop: make(chan struct{}),
		Conf: *cfg,
		DB:   db,
	}

	return app, nil
}

func (p *App) Wait() {
	<-p.stop // 채널에서 값 읽으려고 시도
}

func (p *App) Terminate() {
	close(p.stop)
}

func (p *App) SigHandler() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM) // OS 시그널을 채널로 전달
	go func() {
		sig := <-sigCh
		logger.Info("Arrived terminate signal: ", sig)
		p.Terminate()
	}()
}
