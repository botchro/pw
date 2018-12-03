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
	"errors"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get a registered password with the tag",
	Long: `Get a registered password with the tag.
	The returned password will copied in the clipboard.
	Example:
	$ pw get example.com
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Usage: pw get <TAG>")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		tag := args[0]
		ps := Password{Tag: tag, Password: ""}
		masterPassword, err := ps.GetMasterPassword()
		if err != nil {
			return err
		}

		password, err := ps.GetPassword(masterPassword)
		if err != nil {
			return err
		}

		if err := clipboard.WriteAll(password); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
