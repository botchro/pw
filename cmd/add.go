// Copyright © 2018 Motohiro Nakamura <private.mnakamura@gmail.com>
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
	"fmt"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "register a new password",
	Long: `register a new password.
	ex) pw add example.com password1234
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("Usage: pw add <TAG> <PASSWORD>")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ps := Password{Tag: args[0], Password: args[1]}
		masterPassword, err := ps.GetMasterPassword()
		if err != nil {
			return err
		}
		if err := ps.Register(masterPassword); err != nil {
			return err
		}
		fmt.Println("Successfully registered new password for tag", ps.Tag)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
