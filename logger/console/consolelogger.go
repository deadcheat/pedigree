package console

import "go.uber.org/zap"
import "github.com/deadcheat/pedigree/app"

// ZapLogger Zapを利用してConsoleにロギングする
type ZapLogger struct {
	Logger *zap.Logger
	Tag    string
	Name   string
}

// NewZapLogger
func NewZapLogger(t string, n string) *ZapLogger {
	return &ZapLogger{
		Logger: app.Value.Logger,
		Tag:    t,
		Name:   n,
	}
}

// Log LogInfo with zap
func (z *ZapLogger) Log(o interface{}) (err error) {
	defer z.Logger.Sync()
	z.Logger.Info(z.Tag, zap.Any(z.Name, o))
	return
}
