package app

import (
	"log"
	"os"

	"github.com/fluent/fluent-logger-golang/fluent"
	"go.uber.org/zap"
)

// Env globalレベルのObjectTransportを行うためのもの
var Env *EnvStruct

// ErrLogger 標準エラー出力にエラーを吐くために用意する
var ErrLogger *log.Logger

func init() {
	Env = &EnvStruct{}
	ErrLogger = log.New(os.Stderr, "[ERROR]", log.LstdFlags|log.Lshortfile)
}

// EnvStruct アプリケーション実行においてpackage-global
type EnvStruct struct {
	Logger     *zap.Logger
	Fluent     *fluent.Fluent
	ServerHost *string
	ServerPort *int
	Tag        *string
	ObjectName *string
	FluentHost *string
	FluentPort *int
}

const (
	// TrackingTag トラッキング用のタグ名のデフォルト値
	TrackingTag = "tracking.request"
	// RequestDataname トップレベルの出力オブジェクト名のデフォルト値
	RequestDataname = "RequestData"
)

// EstablishFluent Try to establish connection to fluentd
func EstablishFluent() *fluent.Fluent {
	host := *Env.FluentHost
	port := *Env.FluentPort
	if host == "" || port == 0 {
		return nil
	}
	f, err := fluent.New(fluent.Config{FluentPort: port, FluentHost: host})
	if err != nil {
		ErrLogger.Printf("Error occured when establish connection to fluentd, err : %v \n", err)
		return nil
	}
	return f
}
