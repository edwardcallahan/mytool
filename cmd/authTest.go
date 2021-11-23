/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

// authTestCmd represents the authTest command
var authTestCmd = &cobra.Command{
	Use:   "authTest",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("authTest called")
	},
}

func init() {
	rootCmd.AddCommand(authTestCmd)

}

const (
	authurl = "/cloudstateengine.lightbend.com/v1alpha/users:authenticate"
)

type (
	authConfig struct {
		URL      string `yaml:"url"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	}

	authTestConfig struct {
		Testapps struct {
			AuthTest struct {
				Enabled bool       `yaml:"enabled"`
				Config  authConfig `yaml:"config"`
			}
		}
	}
	config struct {
		Envconfigs map[string]*authTestConfig `yaml:"envs"`
	}
)

// Authenticate function to authenticate
func Authenticate(args interface{}) (string, error) {
	var buff bytes.Buffer
	a, ok := args.(authConfig)
	if !ok {
		return "", fmt.Errorf("%s", "Unable to format errors")
	}

	if a.Password == "" {
		return "", fmt.Errorf("%s", "ERROR password missing for user")
	}

	//url := "https://api.cloudstate.com/cloudstateengine.lightbend.com/v1alpha/users:authenticate"
	url := "https://" + a.URL + authurl
	buff.WriteString(fmt.Sprintf("Authenticate %s to %s", a.Username, url))
	message := map[string]interface{}{
		"password": map[string]string{
			"email_or_friendly_name": a.Username,
			"password":               a.Password,
		},
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		return buff.String(), err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return buff.String(), err
	}

	if resp.StatusCode != 200 {
		return buff.String(), fmt.Errorf("status:%d authentication failed  %s @ %s", resp.StatusCode, a.Username, url)
	}
	//fmt.Println("Severity: INFO, Message: %s ", buff.String())

	return buff.String(), nil
}
