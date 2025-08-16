package app

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abcfe/abcfe-node/api/rest"
	"github.com/abcfe/abcfe-node/common/crypto"
	"github.com/abcfe/abcfe-node/common/logger"
	conf "github.com/abcfe/abcfe-node/config"
	"github.com/abcfe/abcfe-node/core"
	"github.com/abcfe/abcfe-node/storage"
	"github.com/abcfe/abcfe-node/wallet"
	"github.com/syndtr/goleveldb/leveldb"
)

// var port = flag.Int("port", 3000, "Set port of the server")
var mode = flag.String("mode", "rest", "Choose between 'boot' and 'validator'")
var configPath = flag.String("config", "", "Set path of the config file")

type App struct {
	stop       chan struct{}
	Conf       conf.Config
	DB         *leveldb.DB // db내의 mutex는 복사되면 안됨
	BlockChain *core.BlockChain
	restServer *rest.Server // 추가: REST API 서버 필드
	Wallet     *wallet.WalletManager
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
		return nil, err
	}

	wallet, err := wallet.InitWallet(cfg)
	if err != nil {
		logger.Error("Failed to load wallet: ", err)
		return nil, err
	}
	logger.Info("wallet imported: ", crypto.AddressTo0xPrefixString(wallet.Wallet.Accounts[0].Address))

	bc, err := core.NewChainState(db, cfg)
	if err != nil {
		logger.Error("failed to initailze chain state: ", err)
		return nil, err
	}

	app := &App{
		stop:       make(chan struct{}),
		Conf:       *cfg,
		DB:         db,
		BlockChain: bc,
		Wallet:     wallet,
	}

	// REST API 서버 초기화
	app.restServer = rest.NewServer(app.Conf.Server.RestPort, app.BlockChain)

	return app, nil
}

func (p *App) NewRest() error {
	// REST API 서버 시작
	if err := p.restServer.Start(); err != nil {
		return fmt.Errorf("failed to start REST API server: %w", err)
	}

	logger.Info("All services started")
	return nil
}

// Cleanup 애플리케이션 정리
func (p *App) Cleanup() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// REST API 서버 종료
	if p.restServer != nil {
		if err := p.restServer.Stop(ctx); err != nil {
			logger.Error("Error stopping REST API server:", err)
		}
	}

	// DB 연결 닫기
	if p.DB != nil {
		if err := p.DB.Close(); err != nil {
			logger.Error("Error closing DB connection:", err)
		}
	}

	logger.Info("All resources cleaned up")
}

func (p *App) Wait() {
	<-p.stop // 채널에서 값 읽으려고 시도
}

func (p *App) Terminate() {
	p.Cleanup() // 자원 정리 후 종료
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
