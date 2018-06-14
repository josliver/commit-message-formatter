// Copyright © 2018 Rodrigo Navarro <rodrigonavarro23@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type option struct {
	Value string
	Desc  string
}
type keyValue struct {
	Key   string
	Value string
}

var cfgFile string
var variables []keyValue

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cfm",
	Short: "Generate commit message for your repo",
	Long: `CFM (Commit Format Message):
Generate a formated message for your repo using common notations for:
	- Features
	- Fixes
	- Refactoring
	- Tests`,
	PreRun: func(cmd *cobra.Command, args []string) { promptList() },
	Run: func(cmd *cobra.Command, args []string) {
		p := parseTemplate(viper.GetString("template"))
		commit(p)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cfm.yaml)")
}

func loadLocalConfigFile(name string) error {
	viper.AddConfigPath("./")
	viper.SetConfigType("yaml")
	viper.SetConfigName(name)
	return viper.ReadInConfig()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	err := loadLocalConfigFile("default")
	checkErr(err)

	projectDir, err := os.Getwd()
	checkErr(err)
	projectConfigFile := projectDir + "/.cfm.yaml"

	if _, err := os.Stat(projectConfigFile); err == nil {
		cfgFile = projectConfigFile
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		checkErr(err)

		// Search config in home directory with name ".cfm" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cfm")
	}

	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		defaultFlow := viper.GetString("default")
		if defaultFlow != "" {
			err := loadLocalConfigFile(defaultFlow)
			checkErr(err)
		}
	}
}
