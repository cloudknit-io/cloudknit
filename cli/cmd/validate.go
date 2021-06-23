/*
Copyright © 2021 CompuZest <EMAIL ADDRESS>

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
	"errors"
	"fmt"
	"github.com/compuzest/zlifecycle-cli/common"
	v1alpha1 "github.com/compuzest/zlifecycle-cli/types"
	"github.com/compuzest/zlifecycle-cli/validators"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate file",
	Example: "validate path/to/file.yaml",
	Aliases: []string{"v"},
	Args: cobra.ExactArgs(1),
	Short: "Validate a zLifecycle k8s YAML resource",
	Long: `Performs validation on zLifecycle Company, Team and Environment k8s CRD YAMLs.`,
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		dat, err := readFile(path)
		common.HandleError(err, 2)
		kind, err := decodeKind(dat)
		common.HandleError(err, 2)
		err = validate(kind, dat)
		fmt.Println(err)
		common.HandleError(err, 2)
		fmt.Println("Component is valid")
		os.Exit(0)
	},
}

func readFile(path string) ([]byte, error) {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file at %s", path)
		return nil, err
	}
	return dat, nil
}

func decodeKind(dat []byte) (string, error) {
	m := make(map[interface{}]interface{})
	if err := yaml.Unmarshal(dat, &m); err != nil {
		fmt.Println("Error unmarshalling YAML file")
		return "", nil
	}
	return fmt.Sprintf("%v", m["kind"]), nil
}

func validate(kind string, dat []byte) error {
	switch kind {
	case "Environment":
		env := v1alpha1.Environment{}
		if err := yaml.Unmarshal(dat, &env); err != nil {
			fmt.Println("Error while unmarshalling Environment YAML file")
			return err
		}
		if err := validators.ValidateEnvironmentComponents(env.Spec.EnvironmentComponent); err != nil {
			fmt.Println("Error validating Environment resource")
			return err
		}
	default:
		fmt.Printf("Unsupported k8s resource kind: %s", kind)
		return errors.New("unsupported k8s resource")
	}
	return nil
}

func init() {
	rootCmd.AddCommand(validateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// validateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// validateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
