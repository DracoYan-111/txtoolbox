/*
Copyright © 2024 DracoYan-111 <yanlong2944@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"strings"
	config "txtoolbox/cmd/config"
	transaction "txtoolbox/cmd/transaction"
	utils "txtoolbox/cmd/utils"

	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "txToolbox",
	Short: "CLI tool designed for blockchain developers and users",
	Long:  figure.NewFigure("txToolbox", "", true).String(),
	Example: `
utils :Tools
config  :Configuration
transaction :Transaction
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var cfgFile string

func init() {
	// Init function
	cobra.OnInitialize(initConfig)

	// Add command
	rootCmd.AddCommand(utils.UtilsCmd)
	rootCmd.AddCommand(config.ConfigCmd)
	rootCmd.AddCommand(transaction.TransactionCmd)

	// Add flags
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "specify a configuration file (default is: ./.config.env)")

	//Disabling Default Commands
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func initConfig() {
	// Check if config file exists
	if cfgFile != "" {
		if !strings.Contains(cfgFile, ".env") {
			cfgFile = cfgFile + ".env"
		}
	} else {
		cfgFile = ".config.env"
	}

	viper.SetConfigFile(cfgFile)

	// Read in config file
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		for {
			fmt.Println("Create configuration file? (Y/y/N/n)")

			var confirm string
			fmt.Scanln(&confirm)

			switch confirm {
			case "Y", "y":
				viper.WriteConfigAs(cfgFile)
				fmt.Println("Configuration file created.")
				return
			case "N", "n":
				os.Exit(0)
			default:
				fmt.Println("Please enter Y/y or N/n")
			}
		}
	}
}
