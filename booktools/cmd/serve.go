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
	"fmt"

	sv "github.com/TheGrum/booktools/booktools/server"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts booktools as a webservice",
	Long:  `Starts booktools web server.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf(`To access booktools, open a webbrowser and
navigate to http://localhost:%d/%s`, servicePort, "\n\n")
		sv.Listen(processRoot, servicePort)
	},
}

var servicePort int

func init() {
	processCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	serveCmd.Flags().IntVarP(&servicePort, "servicePort", "p", 8080, "Port for webserver")
}
