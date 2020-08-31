//
// Copyright © 2020 alexj@backpocket.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
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
	"log"
	"os"
	"path/filepath"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Display info on authy cmd",
	Long:  `Display info on authy cmd`,

	Run: func(cmd *cobra.Command, args []string) {
		exePath, err := filepath.Abs(os.Args[0])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n\n", exePath)

		fmt.Printf(
			`Application build information
  Build date      : %s
  Build number    : %s
  Git commit      : %s
  Runtime version : %s
  Built on OS     : %s
`, BuildDate, BuildNumber, LatestCommit, RuntimeVer, BuiltOnOs)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
