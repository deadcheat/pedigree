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
	defer app.Value.Logger.Sync()
	app.Value.Logger.Info("tracking.request", zap.Any("request", r))
}
