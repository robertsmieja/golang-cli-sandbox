package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/robertsmieja/golang-cli/util"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

func setup(cmd *cobra.Command, args []string) {
	dockerConfigPath, err := cmd.Flags().GetString(dockerConfigArg)
	if err != nil {
		panic(err)
	}

	dockerConfigBytes, err := ioutil.ReadFile(dockerConfigPath)
	if err, ok := err.(*os.PathError); ok && os.IsNotExist(err) && dockerConfigBytes == nil {
		dockerConfigBytes = make([]byte, 0)
	} else {
		panic(err)
	}

	dockerConfig := new(map[string]interface{})
	if err := json.Unmarshal(dockerConfigBytes, dockerConfig); err != nil {
		if _, ok := err.(*json.SyntaxError); !ok {
			panic(err)
		}
	}

	if util.ValidateDockerDaemonConfig(dockerConfig) {
		return
	}

	util.AddRequiredRegistry(dockerConfig)

	newDockerConfigBytes, err := json.MarshalIndent(dockerConfig, "", "  ")
	if err != nil {
		panic(err)
	}

	stat, err := os.Stat(dockerConfigPath)
	if err != nil {
		panic(err)
	}

	if err = ioutil.WriteFile(dockerConfigPath, newDockerConfigBytes, stat.Mode()); err != nil {
		panic(err)
	}

	fmt.Println("validate called")
}

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: setup,
}

func init() {
	rootCmd.AddCommand(setupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
