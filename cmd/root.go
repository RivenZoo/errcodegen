// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
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
	"github.com/RivenZoo/errcodegen/pkg/config"
	"github.com/RivenZoo/errcodegen/pkg/errcode_def"
	"github.com/RivenZoo/errcodegen/pkg/gen"
	"github.com/RivenZoo/errcodegen/pkg/log"
	"github.com/RivenZoo/errcodegen/pkg/patch"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var commonConfig config.ErrCodeCommonConfig
var errorCodeDefineFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "errcodegen",
	Short: "A code generate tool for generating error code.",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if cfgFile != "" {
			config.MustInit(cfgFile)
		}
		cfg := config.GetConfig()
		err := patchCommonConfig(cfg)
		if err != nil {
			log.Error("patchCommonConfig error %v", err)
			os.Exit(-1)
		}
		if errorCodeDefineFile != "" {
			func() {
				f, err := os.Open(errorCodeDefineFile)
				if err != nil {
					log.Error("error code define file %s not exists", errorCodeDefineFile)
					os.Exit(-1)
				}
				defer f.Close()

				p := errcode_def.NewModuleErrorCodeParser()
				modules, err := p.ParseErrorCodeDefinition(f)
				if err != nil {
					log.Error("parse error code define file error %v", err)
					os.Exit(-1)
				}
				cfg.Modules = append(cfg.Modules, modules...)
			}()
		}

		config.PatchConfig(cfg)
		config.CheckConfig(cfg)
		g := gen.NewCodeGenerator()
		for i := range cfg.Modules {
			err := g.GenerateModuleErrorCode(&cfg.ErrCodeCommonConfig, &cfg.Modules[i])
			if err != nil {
				os.Exit(-1)
			}
		}
	},
}

func patchCommonConfig(cfg *config.ErrCodeConfig) error {
	return patch.CoverStructsField(commonConfig, &cfg.ErrCodeCommonConfig)
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

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.errcodegen.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringVar(&commonConfig.PkgName, "pkg", "errorcodes", "Generated module default package name")
	rootCmd.Flags().StringVar(&commonConfig.AppCode, "appcode", "100", "App error code, three digits")
	rootCmd.Flags().StringVar(&commonConfig.NewErrorFuncPkg, "err_func_pkg", "errors", "New error function package import path")
	rootCmd.Flags().StringVar(&commonConfig.NewErrorFunc, "err_func", "New", "New error function name")
	rootCmd.Flags().StringVar(&errorCodeDefineFile, "err_def", "", "Error code definiton file")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".errcodegen" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".errcodegen")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
