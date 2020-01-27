package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/robertsmieja/golang-cli/util"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

const dockerConfigArg = "docker-config"

func validate(cmd *cobra.Command, _ []string) {
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

	fmt.Println("Config should be:")
	fmt.Println(string(newDockerConfigBytes))
	fmt.Println("validate called")
}

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:     "validate",
	Short:   "Validate current configuration",
	Long:    `Validate current configuration and report issues found`,
	Run:     validate,
	Aliases: []string{"verify"},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	validateCmd.Flags().String(dockerConfigArg, "/etc/docker/daemon.json", "Path to Docker Daemon configuration")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// validateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// validateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
