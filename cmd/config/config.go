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
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ConfigCmd represents the config/config command
var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Set up a configuration file for use",
	Long:  figure.NewFigure("config", "", true).String(),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("config/config called")
		for k, v := range GetConfig() {
			fmt.Printf("%s=%v\n", k, v)
		}
		return nil
	},
}

// ConfigGetCmd represents the config/get command
var ConfigGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the specified configuration file",
	Example: `
config get -k:Get the specified configuration file`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("config/getConfig called")
		value, err := GetConfigByKey(key)
		if err != nil {
			return err
		}
		fmt.Println("K:[", key, "] V:[", value, "]")
		return err
	},
}

// ConfigGetCmd represents the config/add command
var ConfigAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add the specified configuration file",
	Example: `
config add -k -v:Add the specified configuration file`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("config/addConfig called")
		err := AddConfig(key, value)
		if err == nil {
			newValue, errs := GetConfigByKey(key)
			if errs != nil {
				return errs
			}
			fmt.Println("Key:[", key, "] Value:[", newValue, "]")
		}
		return err
	},
}

// ConfigSetCmd represents the config/set command
var ConfigSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set the specified configuration file",
	Example: `
config set -k -v:Set the specified configuration file`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("config/setConfig called")
		err := SetConfigByKey(key, value)
		if err == nil {
			newValue, errs := GetConfigByKey(key)
			if errs != nil {
				return errs
			}
			fmt.Println("Key:[", key, "] Value:[", newValue, "]")
		}
		return err
	},
}

// ConfigDelCmd represents the config/del command
var ConfigDelCmd = &cobra.Command{
	Use:   "del",
	Short: "Delete the specified configuration file",
	Example: `
config del -k:Delete the specified configuration file`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("config/delConfig called")
		err := DelConfigByKey(key)
		return err
	},
}

var key string
var value string

func init() {
	// Add command
	ConfigCmd.AddCommand(ConfigGetCmd)
	ConfigCmd.AddCommand(ConfigSetCmd)
	ConfigCmd.AddCommand(ConfigDelCmd)
	ConfigCmd.AddCommand(ConfigAddCmd)

	// Add flags
	ConfigGetCmd.Flags().StringVarP(&key, "key", "k", "", "key")
	ConfigGetCmd.MarkFlagRequired("key")

	ConfigAddCmd.Flags().StringVarP(&key, "key", "k", "", "key")
	ConfigAddCmd.Flags().StringVarP(&value, "value", "v", "", "value")
	ConfigAddCmd.MarkFlagRequired("key")
	ConfigAddCmd.MarkFlagRequired("value")

	ConfigSetCmd.Flags().StringVarP(&key, "key", "k", "", "key")
	ConfigSetCmd.Flags().StringVarP(&value, "value", "v", "", "value")
	ConfigSetCmd.MarkFlagRequired("key")
	ConfigSetCmd.MarkFlagRequired("value")

	ConfigDelCmd.Flags().StringVarP(&key, "key", "k", "", "key")
	ConfigDelCmd.MarkFlagRequired("key")

}

// Get all config
func GetConfig() map[string]any {
	return viper.AllSettings()
}

// Get config by key
func GetConfigByKey(key string) (string, error) {
	value := viper.GetString(key)
	if value == "" {
		return "", errors.New("key does not exist")
	}
	return value, nil
}

// Add config
func AddConfig(key, value string) error {

	viper.Set(key, value)
	err := viper.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}

// Set config by key
func SetConfigByKey(key, value string) error {
	if viper.Get(key) == nil {
		return fmt.Errorf("key %s already exists", key)
	}

	viper.Set(key, value)
	err := viper.WriteConfig()
	if err != nil {
		return err
	}

	return nil
}

// Delete config by key
// Deleting the configuration file will cause the KEY content of the configuration file to become lowercase.
// However, viper ignores case when reading configuration files, so it will not affect
func DelConfigByKey(key string) error {
	if viper.Get(key) == nil {
		return errors.New("key does not exist")
	}
	// Create a map and delete the specified key
	newAllConfig := make(map[string]any, len(viper.AllKeys()))

	for k, v := range viper.AllSettings() {
		newAllConfig[k] = v
	}

	// Delete the specified key
	delete(newAllConfig, strings.ToLower(key))

	// Create a new file to overwrite the original file
	configFileUsed := viper.ConfigFileUsed()
	files, _ := os.Create(configFileUsed)

	// Defer close the file
	defer files.Close()
	// Reset viper config
	viper.Reset()

	// Write the new config
	for k, v := range newAllConfig {
		viper.Set(k, v)
		files.WriteString(k + "=" + v.(string) + "\n")
	}

	//Set the config file address
	viper.SetConfigFile(configFileUsed)

	return viper.WriteConfig()
}
