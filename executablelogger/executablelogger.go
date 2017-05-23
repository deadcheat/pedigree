package executablelogger

import "github.com/deadcheat/pedigree/chainer"
import "github.com/deadcheat/pedigree/logger"

// ExecutableLogger ロガー実装struct
type ExecutableLogger struct {
	Logger logger.Loggable
}

// NewExecutableLogger ログ実行
func NewExecutableLogger(l logger.Loggable) chainer.Executable {
	return &ExecutableLogger{
		Logger: l,
	}
}

// Execute ログ出力を実行してNextをコールする
func (e *ExecutableLogger) Execute(c chainer.Chainable, o interface{}) error {
	return e.Logger.Log(o)
}
