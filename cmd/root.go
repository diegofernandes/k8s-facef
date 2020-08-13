/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var (
	cfgFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "k8s-facef",
	Short: "k8s-facef is k8s sample application",
	Long:  `k8s-facef is k8s sample application`,

	Run: func(cmd *cobra.Command, args []string) {

		var basicUser = viper.GetString("auth_user")
		var basicPass = viper.GetString("auth_pass")

		addr := fmt.Sprintf(":%v", viper.GetInt("port"))
		log.Printf("Starting hello api on %s", addr)
		if basicUser != "" {
			log.Printf("Using HTTP basic auth ")
		}

		mux := http.NewServeMux()

		mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {

			w.WriteHeader(200)
			w.Write([]byte("ok"))

		})
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

			if basicUser != "" {
				if username, password, ok := r.BasicAuth(); !ok || username != basicUser || password != basicPass {
					log.Printf("User Unauthorized! Got username=%s password=%s Should be username=%s password=%s", username, password, basicUser, basicPass)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
			}

			for i := 0; i < 999999999; i++ {
			}
			runtime.Gosched()
			msg := fmt.Sprintf(viper.GetString("message"), time.Now().Format(viper.GetString("time_layout")))

			fmt.Fprint(w, msg)
		})
		log.Println(http.ListenAndServe(addr, mux))

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./k8s-facef.yaml)")

	rootCmd.Flags().Int("port", 8080, "HTTP listener port")
	rootCmd.Flags().String("message", "Hello cruel world, It is %s", "Welcome message format")
	rootCmd.Flags().String("time_layout", time.UnixDate, "time layout used in the message")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().String("auth_user", "", "HTTP basic auth username")
	rootCmd.Flags().String("auth_pass", "", "HTTP basic auth password")

	viper.BindPFlags(rootCmd.PersistentFlags())

	viper.BindPFlags(rootCmd.Flags())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// directory path of config file
		viper.AddConfigPath("/etc/k8s-facef")
		// directory path of config file
		viper.AddConfigPath(".")
		// name of config file (without extension)
		viper.SetConfigName("k8s-facef")

	}

	// read in environment variables that match
	viper.AutomaticEnv()
	// all environment variables that contains _ will be replaced by - (example - env: HOST_NAME -> access: host-name)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		log.Println(err)
	}

}
