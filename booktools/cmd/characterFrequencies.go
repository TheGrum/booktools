// Copyright © 2018 Howard C. Shaw III <howardcshaw@gmail.com>
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
	bt "github.com/TheGrum/booktools"

	"github.com/spf13/cobra"
)

// characterFrequenciesCmd represents the characterFrequencies command
var characterFrequenciesCmd = &cobra.Command{
	Use:   "characterFrequencies",
	Short: "Lists the characters and the frequency with which they appear.",
	Run: func(cmd *cobra.Command, args []string) {
		bt.PrintNameFrequency(bt.CharacterFrequencies(processRoot, 3, 1))
	},
}

func init() {
	processCmd.AddCommand(characterFrequenciesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// characterFrequenciesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// characterFrequenciesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
