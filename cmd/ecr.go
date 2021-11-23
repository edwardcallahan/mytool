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
	"log"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

type EcrReturn struct {
	ImageIds []struct {
		ImageDigest string `json:"imageDigest"`
		ImageTag    string `json:"imageTag,omitempty"`
	} `json:"imageIds"`
}

var ecrCmd = &cobra.Command{
	Use:   "ecr",
	Short: "Find the digest for the given ECR image tag.",
	Long: `Given AWS credentials and an image tag, search for matching tags and return digest if found.

	Requires:  --repo, and --image 
	
	Defaults:  --tag latest and --region us-east-1`,
	Run: func(cmd *cobra.Command, args []string) {
		if image != "image" && repo != "myECRrepo" {
			fmt.Println("ECR Tag Finder, seeking the digest for:", image, "with tag:", tag, "in", repo)
		} else {
			// flag.Usage() causing problems so...
			fmt.Println("Please see help using `ecr --help`.")
			return
		}

		ecr := exec.Command("aws", "ecr", "list-images", "--repository-name", repo, "--region", region)

		var stdout, stderr bytes.Buffer
		ecr.Stdout = &stdout
		ecr.Stderr = &stderr
		err := ecr.Run()
		listingStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())

		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n %s\n", err, errStr)
			return
		}

		ecrRet := EcrReturn{}
		buffer := []byte(listingStr)

		if err := json.Unmarshal(buffer, &ecrRet); err != nil {
			fmt.Printf("error unmarshaling image JSON: %v\n", err)
			return
		}

		var images = ecrRet.ImageIds
		 for k, v := range images{
			key := fmt.Sprintf("%#v", k)
			value := fmt.Sprintf("%#v", v)
			thisTag := string(value)
			if (strings.Contains(thisTag,tag)) {
				fmt.Println("Found digest for tag:", tag)
				fmt.Println("***",key, value)
			 }
		 }

	},
}

func extract(b []byte) string {
	var m map[string]string
	err := json.Unmarshal(b ,&m)
	if err!= nil {
		fmt.Println("error unmarshalling", err)
	}
	fmt.Println()
	return "hi"
}

var image, tag, region, repo string

func init() {
	rootCmd.AddCommand(ecrCmd)
	ecrCmd.Flags().StringVar(&image, "image", "image", "--image=<image name>")
	ecrCmd.Flags().StringVar(&tag, "tag", "latest", "--tag=latest")
	ecrCmd.Flags().StringVar(&repo, "repo", "myECRrepo", "--repo=<repo name>")
	ecrCmd.Flags().StringVar(&region, "region", "us-east-1", "--region=<aws region>")
}
