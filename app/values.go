package app

import "go.uber.org/zap"

// AppEnv アプリケーション実行においてpackage-global
type AppEnv struct {
	Logger *zap.Logger
}
