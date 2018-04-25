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
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	bt "github.com/TheGrum/booktools"

	"github.com/spf13/cobra"
)

// processCmd represents the process command
var processCmd = &cobra.Command{
	Use:   "process",
	Short: "Process the specified file",
	Long: `Reads the specified file, tokenizes and chunks it
in preparation for further procssing.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			Process(os.Stdin)
		}
		for _, arg := range args {
			if arg == "-" {
				processRoot = Process(os.Stdin)
			} else {
				file, err := os.Open(arg)
				if err != nil {
					log.Fatalf("Error opening file to process: %v", err)
				}
				processRoot = Process(file)
			}
		}
	},
}

var processRoot *bt.Chunk
var chapterRegex string

func init() {
	rootCmd.AddCommand(processCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//processCmd.PersistentFlags().String("file", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// processCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	processCmd.PersistentFlags().StringVarP(&chapterRegex, "chapterRegex", "r", "", "Regular expression which if matched on a line will trigger a chapter.")
}

func Process(input io.Reader) *bt.Chunk {
	data, err := ioutil.ReadAll(input)
	if err != nil {
		log.Fatal(err)
	}
	/*
		file, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	*/

	chunks := make(chan *bt.Chunk, 10)
	out := make(chan *bt.Chunk)

	var reg *regexp.Regexp
	if chapterRegex != "" {
		reg, err = regexp.Compile(chapterRegex)
		if err != nil {
			log.Fatalf("Failed to compile chapter-matching regular expression [%v]", chapterRegex)
		}
	}

	chunker := bt.NewChunker(strings.NewReader(string(data)), chunks)
	chunker.OnBeforeSentence = func(c *bt.Chunker, s string) {
		if !(reg == nil) {
			if reg.Match([]byte(s)) {
				c.Chapter()
			}
		} else {
			// Here, define what counts as a new Chapter
			if strings.HasPrefix(strings.TrimSpace(s), "Chapter ") {
				c.Chapter()
			} else if strings.Contains(s, "# Part") || strings.Contains(s, "#Part") {
				c.Chapter()
			} else if strings.HasSuffix(s, "2011") {
				c.Chapter()
			}
		}
	}
	scanner := bufio.NewScanner(chunker)
	go bt.DigestChunks(chunks, out)
	for scanner.Scan() {
		//fmt.Println("[" + scanner.Text() + "]")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	//fmt.Println("Fetching root")
	root := <-out
	return root
}
