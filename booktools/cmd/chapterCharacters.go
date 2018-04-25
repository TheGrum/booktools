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
	"github.com/spf13/cobra"
)

// chapterCharactersCmd represents the chapterCharacters command
var chapterCharactersCmd = &cobra.Command{
	Use:   "chapterCharacters",
	Short: "Lists the characters in each chapter",
	Long:  `Lists the top appearing characters per chapter`,
	Run: func(cmd *cobra.Command, args []string) {
		processRoot.PrintTopXCharactersPerChapter(includeSentences, includeXthSentence, topX, tabDelimit, wordCount)
	},
}

var topX int
var includeXthSentence int
var tabDelimit bool
var wordCount bool

func init() {
	processCmd.AddCommand(chapterCharactersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chapterCharactersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// chapterCharactersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	chapterCharactersCmd.Flags().BoolVarP(&includeSentences, "includeSentences", "s", false, "Include first sentence of each chapter")
	chapterCharactersCmd.Flags().BoolVarP(&tabDelimit, "tabDelimit", "b", false, "Tab delimit the output")
	chapterCharactersCmd.Flags().BoolVarP(&wordCount, "wordCount", "w", false, "Include word count of each chapter")
	chapterCharactersCmd.Flags().IntVarP(&includeXthSentence, "includeXth", "x", -1, "Include Xth sentence of each chapter")
	chapterCharactersCmd.Flags().IntVarP(&topX, "topX", "t", 3, "Show top X characters")
}
