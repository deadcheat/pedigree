package app

import "go.uber.org/zap"

var Value *AppEnv

// AppEnv アプリケーション実行においてpackage-global
type AppEnv struct {
	Logger     *zap.Logger
	ServerHost *string
	ServerPort *int
}
