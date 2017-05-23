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

	"github.com/deadcheat/pedigree/app"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// loggerCmd represents the logger command
var loggerCmd = &cobra.Command{
	Use:   "logger",
	Short: "startup request-logging server",
	Long:  `startup request-logging server`,
	Run:   startLogging,
}

func init() {
	RootCmd.AddCommand(loggerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loggerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loggerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	app.Value = &app.AppEnv{}
	app.Value.Logger, _ = zap.NewDevelopment()
	app.Value.ServerHost = loggerCmd.Flags().StringP("host", "H", "localhost", "specify hostname, default: localhost")
	app.Value.ServerPort = loggerCmd.Flags().IntP("port", "p", 3000, "specify portnum, default: 3000")
}

func startLogging(cmd *cobra.Command, args []string) {
	hostName := fmt.Sprintf("%s:%d", *app.Value.ServerHost, *app.Value.ServerPort)
	log.Printf("server start in %s \n", hostName)
	http.HandleFunc("/", loggingHandler)
	if err := http.ListenAndServe(
		hostName,
		nil); err != nil {
		defer app.Value.Logger.Sync()
		app.Value.Logger.Error("http-error occured", zap.Error(err))
	}
}

func loggingHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	defer app.Value.Logger.Sync()
	// Body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

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
	// log through zap
	app.Value.Logger.Info("tracking.request",
		zap.Any("RequestData", o.Data),
	)
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
