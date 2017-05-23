package logger

import "github.com/deadcheat/pedigree/chainer"

// Loggable ログInterface
type Loggable interface {
	Log(o interface{}) error
}

// ExecutableLogger ロガー実装struct
type ExecutableLogger struct {
	Logger Loggable
}

// NewExecutableLogger ログ実行
func NewExecutableLogger(l Loggable) chainer.Executable {
	return &ExecutableLogger{
		Logger: l,
	}
}

// Execute ログ出力を実行してNextをコールする
func (e *ExecutableLogger) Execute(c chainer.Chainable, o interface{}) error {
	return e.Logger.Log(o)
}
