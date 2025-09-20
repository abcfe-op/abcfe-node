package logger

import (
	"bytes"
	"encoding/json"
	"net/http"

	// "encoding/json"
	"fmt"
	"os"

	"time"

	conf "github.com/abcfe-op/abcfe-node/config"

	// "github.com/gin-gonic/gin"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger
var stag string
var cf *conf.Config

func InitLogger(cfg *conf.Config) error {
	now := time.Now()
	lPath := fmt.Sprintf("%s_%s.log", cfg.LogInfo.Fpath, now.Format("2006-01-02"))
	cf = cfg

	rotator, err := rotatelogs.New(
		lPath,
		rotatelogs.WithMaxAge(time.Duration(cfg.LogInfo.MaxAgeHour)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(cfg.LogInfo.RotateHour)*time.Hour))
	if err != nil {
		return err
	}

	encCfg := zapcore.EncoderConfig{
		TimeKey:        "date",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	w := zapcore.AddSync(rotator)
	cw := zapcore.AddSync(os.Stdout)
	var core zapcore.Core
	// stag = cfg.Info.RunMode
	stag = cfg.Common.Mode
	if stag == "alpha" {
		core = zapcore.NewTee(
			zapcore.NewCore(zapcore.NewJSONEncoder(encCfg), w, zap.DebugLevel),
			zapcore.NewCore(zapcore.NewConsoleEncoder(encCfg), cw, zap.DebugLevel),
		)
	} else {
		core = zapcore.NewCore(zapcore.NewJSONEncoder(encCfg), w, zap.InfoLevel)
	}
	logger = zap.New(core)

	logger.Info("logging init file start")
	return nil
}

// msg := fmt.Sprint(ctx...)
// switch level {
// case zapcore.DebugLevel:
// 	logger.Debug("debug", zap.String("Debug", msg))

func Debug(ctx ...interface{}) {
	var b bytes.Buffer
	for _, str := range ctx {
		b.WriteString(fmt.Sprintf("%v", str))
	}

	logger.Debug("debug", zap.String("Debug", b.String()))
}

// Info is a convenient alias for Root().Info
func Info(ctx ...interface{}) {
	var b bytes.Buffer

	for _, str := range ctx {
		b.WriteString(fmt.Sprintf("%v", str))
	}
	// logger.Info("info", zap.String("Info", b.String()))
	logger.Info("info", zap.String("Info", b.String()))
}

// Warn is a convenient alias for Root().Warn
func Warn(ctx ...interface{}) {
	var b bytes.Buffer
	for _, str := range ctx {
		b.WriteString(fmt.Sprintf("%v", str))
	}

	logger.Warn("warn", zap.String("Warn", b.String()))
}

// Error is a convenient alias for Root().Error
func Error(ctx ...interface{}) {
	var b bytes.Buffer
	for _, str := range ctx {
		b.WriteString(fmt.Sprintf("%v", str))
	}

	logger.Error("error", zap.String("Err", b.String()))
	if stag != "alpha" {
		go sendTelegramAlert(cf, b.String())
	}
}

func Crit(ctx ...interface{}) {
	var b bytes.Buffer
	for _, str := range ctx {
		b.WriteString(fmt.Sprintf("%v", str))
	}

	logger.Fatal("panic", zap.String("Crit", b.String()))
	if stag != "alpha" {
		go sendTelegramAlert(cf, b.String())
	}
}

func sendTelegramAlert(cf *conf.Config, body string) bool {
	path, _ := os.Getwd()
	var msg string
	telKey := cf.LogInfo.DevTelKey
	chatId := cf.LogInfo.DevChatId

	if cf.Common.Mode == "prod" {
		msg = "[" + cf.Common.ServiceName + "_" + cf.Common.Mode + "] " + "!!!From Prod-live stage!!! : \n" + body + "\nModule : " + path
		telKey = cf.LogInfo.ProdTelKey
		chatId = cf.LogInfo.ProdChatId
	} else if cf.Common.Mode == "dev" {
		msg = "[" + cf.Common.ServiceName + "_" + cf.Common.Mode + "] " + "From dev stage : \n" + body + "\nModule : " + path
	} else {
		msg = "[" + cf.Common.ServiceName + "_" + cf.Common.Mode + "] " + " Message : \n" + body + "\nModule : " + path
	}

	pbytes, _ := json.Marshal(map[string]interface{}{"chat_id": chatId, "text": msg})
	buff := bytes.NewBuffer(pbytes)
	if _, err := http.Post(telKey, "application/json", buff); err != nil {
		return false
	}

	return true
}

// 에러 핸들링
func HandleErr(err error) {
	if err != nil {
		Error(err)
	}
}
