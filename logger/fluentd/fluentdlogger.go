package fluentd

import (
	"github.com/deadcheat/pedigree/app"
	"github.com/fluent/fluent-logger-golang/fluent"
)

// FluentLogger fluent-logger-golang/fluentをアダプターするためのもの
type FluentLogger struct {
	Logger *fluent.Fluent
	Tag    string
	Name   string
}

// NewFluentlogger FluentLoggerを新規に生成
func NewFluentlogger(t string, n string) *FluentLogger {
	if app.Env.Fluent == nil {
		return nil
	}
	return &FluentLogger{
		Logger: app.Env.Fluent,
		Tag:    t,
		Name:   n,
	}
}

// Log implemnt of Loggable fluentdへのログ出力を行う
func (f *FluentLogger) Log(o interface{}) (err error) {
	if f == nil {
		return
	}
	m := map[string]interface{}{
		f.Name: o,
	}
	if err := f.Logger.Post(f.Tag, m); err != nil {
		return err
	}
	return
}
