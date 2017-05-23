package logger

// Loggable ログInterface
type Loggable interface {
	Log(o interface{}) error
}
