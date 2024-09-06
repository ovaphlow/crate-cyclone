package infrastructure

import (
	"log"
	"log/slog"
	"os"
)

var Slogger *slog.Logger

func InitSlog() {
	Slogger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(Slogger)
	log.Println("日志组件初始化完成")
}
