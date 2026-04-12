package server

import (
	"ElainaWeb/global"

	"go.uber.org/zap"
)

type server interface {
	ListenAndServe() error
}

func RunServer() {
	addr := global.Config.System.Addr()
	Router := RouterInit()

	s := initServer(addr, Router)
	global.Logger.Info("服务启动成功", zap.String("address", addr))
	if err := s.ListenAndServe(); err != nil {
		global.Logger.Error("服务异常停止", zap.Error(err))
	}
}
