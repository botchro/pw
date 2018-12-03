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

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "show the list of tags",
	Long: `Show the list of tags.
	Example:
	$ pw list
	example.com
	foo.co.jp
	bar.io
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ps := new(Password)
		masterPassword, err := ps.GetMasterPassword()
		if err != nil {
			return err
		}

		tags, err := ps.GetTagList(masterPassword)
		if err != nil {
			return err
		}

		for _, tag := range tags {
			fmt.Println(tag)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
