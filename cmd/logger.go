// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/deadcheat/pedigree/actionstore"
	"github.com/deadcheat/pedigree/app"
	"github.com/deadcheat/pedigree/executablelogger"
	"github.com/deadcheat/pedigree/logger/console"
	"github.com/deadcheat/pedigree/logger/fluentd"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// loggerCmd represents the logger command
var loggerCmd = &cobra.Command{
	Use:   "logger",
	Short: "startup request-logging server",
	Long: `
__________             .___.__                                           .____                                               ._.
\______   \  ____    __| _/|__|   ____  _______   ____    ____           |    |      ____     ____     ____    ____  _______ | |
 |     ___/_/ __ \  / __ | |  |  / ___\ \_  __ \_/ __ \ _/ __ \   ______ |    |     /  _ \   / ___\   / ___\ _/ __ \ \_  __ \| |
 |    |    \  ___/ / /_/ | |  | / /_/  > |  | \/\  ___/ \  ___/  /_____/ |    |___ (  <_> ) / /_/  > / /_/  >\  ___/  |  | \/ \|
 |____|     \___  >\____ | |__| \___  /  |__|    \___  > \___  >         |_______ \ \____/  \___  /  \___  /  \___  > |__|    __
                \/      \/     /_____/               \/      \/                  \/        /_____/  /_____/       \/          \/
pedigree-logger is one & only function of pedigree.
startup request-logging server.`,
	Run: startLogging,
}

func init() {
	RootCmd.AddCommand(loggerCmd)
	// Produtionのconfigそのままでいいと思うんだけどOutputPathsがstderrなので標準エラーに吐かれちゃうのだけ変更したい
	conf := zap.NewProductionConfig()
	conf.OutputPaths = []string{"stdout"}
	app.Env.Logger, _ = conf.Build()
	app.Env.ServerHost = loggerCmd.Flags().StringP("host", "H", "localhost", "Specify hostname, default: localhost")
	app.Env.ServerPort = loggerCmd.Flags().IntP("port", "p", 3000, "Specify portnum, default: 3000")
	app.Env.Tag = loggerCmd.Flags().StringP("tag", "t", app.TrackingTag, "Tag name that should be passed to fluentd, default: "+app.TrackingTag)
	app.Env.ObjectName = loggerCmd.Flags().StringP("name", "n", app.RequestDataname, "Top-Level object's name that will be logged, default: "+app.RequestDataname)
	app.Env.FluentHost = loggerCmd.Flags().String("fluent-host", "", "Specify fluentd host default is not set and never access fluentd")
	app.Env.FluentPort = loggerCmd.Flags().Int("fluent-port", 0, "Specify fluentd port default is not set and never access fluentd")
	app.Env.CORSConfFile = loggerCmd.Flags().StringP("cors-config", "c", "", "Specify config file path")

}

func corsEnv(path *string) (env *cors.Cors) {
	// この辺別途funcに切り出すべき
	if path == nil || *path == "" {
		// 未指定の場合は設定なしで終わり
		return
	}
	corsPath, err := filepath.Abs(*path)
	if err != nil {
		// initでのエラーはFatal呼んだほうが良い気がしている
		log.Fatalf("Error occured when reading cors-config file %s. errors: %v \n", corsPath, err)
	} else {
		var conf app.CORSEnv
		_, err := toml.DecodeFile(corsPath, &conf)
		if err != nil {
			log.Fatalf("Error occured when reading cors-config file %s. errors: %v \n", corsPath, err)
		}
		env = cors.New(cors.Options{
			AllowedOrigins:     conf.AllowedOrigins,
			AllowedMethods:     conf.AllowedMethods,
			AllowedHeaders:     conf.AllowedHeaders,
			ExposedHeaders:     conf.ExposedHeaders,
			AllowCredentials:   conf.AllowCredentials,
			MaxAge:             conf.MaxAge,
			OptionsPassthrough: conf.OptionsPassthrough,
			Debug:              conf.Debug,
		})
	}
	return
}

func startLogging(cmd *cobra.Command, args []string) {
	// config load
	corsEnv := corsEnv(app.Env.CORSConfFile)

	// establish fluent connection
	app.Env.Fluent = app.EstablishFluent()
	if app.Env.Fluent != nil {
		defer app.Env.Fluent.Close()
	}
	hostName := fmt.Sprintf("%s:%d", *app.Env.ServerHost, *app.Env.ServerPort)
	log.Printf("server start in %s \n", hostName)
	mux := http.NewServeMux()
	mux.HandleFunc("/", loggingHandler)
	var ha http.Handler
	if corsEnv != nil {
		app.Env.CORSEnabled = true
		ha = corsEnv.Handler(mux)
	} else {
		ha = mux
	}
	if err := http.ListenAndServe(hostName, ha); err != nil {
		defer app.Env.Logger.Sync()
		app.Env.Logger.Error("http-error occured", zap.Error(err))
	}
}

func loggingHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	defer app.Env.Logger.Sync()
	allowOrigin := w.Header().Get("Access-Control-Allow-Origin")
	if app.Env.CORSEnabled && allowOrigin == "" {
		// allowがなければ終わり
		return
	}
	// Body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	go func(r *http.Request, b []byte) {
		// Form
		form := NewKeyValueArray()
		for k, v := range r.Form {
			form.Add(map[string]interface{}{
				k: v,
			})
		}
		post := NewKeyValueArray()
		for k, v := range r.PostForm {
			post.Add(map[string]interface{}{
				k: v,
			})
		}

		// Header
		header := NewKeyValueArray()
		for k, v := range r.Header {
			header.Add(map[string]interface{}{
				k: v,
			})
		}

		// Cookies
		cookies := NewKeyValueArray()
		requestedCookies := r.Cookies()
		for i := range requestedCookies {
			c := requestedCookies[i]
			cookies.Add(map[string]interface{}{
				c.Domain: c.String(),
			})
		}

		// Request
		request := map[string]interface{}{
			"Method":     r.Method,
			"URL":        r.URL.String(),
			"Host":       r.Host,
			"Proto":      r.Proto,
			"RequestURI": r.RequestURI,
			"RemoteAddr": r.RemoteAddr,
			"Referer":    r.Referer(),
			"UA":         r.UserAgent(),
		}

		o := NewKeyValueArray()
		o.Add(map[string]interface{}{"Request": request})
		o.Add(map[string]interface{}{"Header": header.Data})
		o.Add(map[string]interface{}{"Body": string(b)})
		o.Add(map[string]interface{}{"Form": form.Data})
		o.Add(map[string]interface{}{"PostForm": post.Data})
		o.Add(map[string]interface{}{"Cookie": cookies.Data})

		as := actionstore.NewActionStore()
		as.Object = o.Data
		as.Add(executablelogger.NewExecutableLogger(console.NewZapLogger(
			*app.Env.Tag, *app.Env.ObjectName,
		)))
		as.Add(executablelogger.NewExecutableLogger(fluentd.NewFluentlogger(
			*app.Env.Tag, *app.Env.ObjectName,
		)))
		if err := as.Next(); err != nil {
			app.ErrLogger.Printf("Error occured in parallel routine, err: %v \n", err)
		}
	}(r, b)
}

// KVArray Key-Value形式のペアオブジェクトを格納する
type KVArray struct {
	Data []map[string]interface{} `json:"data,omitdata"`
}

// NewKeyValueArray KVArray生成
func NewKeyValueArray() *KVArray {
	return &KVArray{
		Data: make([]map[string]interface{}, 0),
	}
}

// Add 要素を追加する
func (k *KVArray) Add(kv map[string]interface{}) {
	k.Data = append(k.Data, kv)
	return
}
