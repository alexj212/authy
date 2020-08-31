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
	"os"
)

// wipeCmd represents the wipe command
var wipeCmd = &cobra.Command{
	Use:   "wipe",
	Short: "remove ~/.authy*.json cache files",
	Long: `Remove authy cache files
.authy.json
.authycache.json

`,
	Run: func(cmd *cobra.Command, args []string) {
		wipeExec(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(wipeCmd)
}

func wipeExec(cmd *cobra.Command, args []string) {

	fpath, err := ConfigPath(configFileName)
	if err != nil {
		fmt.Printf("unable to get config file: %s error: %v\n", configFileName, err)
		return
	}

	if fileExists(fpath) {
		err = os.Remove(fpath)
		if err != nil {
			fmt.Printf("unable to delete config file: %s error: %v\n", fpath, err)
			return
		}

		fmt.Printf("deleted config file: %s\n", fpath)
	}

	fpath, err = ConfigPath(cacheFileName)
	if err != nil {
		fmt.Printf("unable to get cache file: %s error: %v\n", cacheFileName, err)
		return
	}

	if fileExists(fpath) {
		err = os.Remove(fpath)
		if err != nil {
			fmt.Printf("unable to delete cache file: %s error: %v\n", fpath, err)
			return
		}
		fmt.Printf("deleted cache file: %s\n", fpath)
	}

	fmt.Printf("Removed config & cache files\n")
}
