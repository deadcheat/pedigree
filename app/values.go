package app

import (
	"log"

	"github.com/fluent/fluent-logger-golang/fluent"
	"go.uber.org/zap"
)

// Value globalレベルのObjectTransportを行うためのもの
var Value *AppEnv

// AppEnv アプリケーション実行においてpackage-global
type AppEnv struct {
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
	host := *Value.FluentHost
	port := *Value.FluentPort
	if host == "" || port == 0 {
		return nil
	}
	f, err := fluent.New(fluent.Config{FluentPort: port, FluentHost: host})
	if err != nil {
		log.Printf("Error occured when establish connection to fluentd, err : %v \n", err)
		return nil
	}
	return f
}
