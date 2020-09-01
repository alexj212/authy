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
	"log"

	"github.com/spf13/cobra"
)

// refreshCmd represents the refresh command
var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh token cache",
	Long: `When you add a new token.

You can use this cmd to refresh local token cache`,
	Run: func(cmd *cobra.Command, args []string) {
		devInfo, _, err := Initialize()
		if err != nil {
			log.Fatal("Load device info failed", err)
		}

		getTokensFromAuthyServer(devInfo)
	},
}

func init() {
	rootCmd.AddCommand(refreshCmd)
}
