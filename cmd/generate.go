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
	//"time"

	"github.com/spf13/cobra"
)

// ./bin/authy generate PalVPN "/home/alexj/bin/vpnup.sh [AUTHCODE]"

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate [TokenName]",
	Short: "generate a otp token",
	Long: `generate a otp token
`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			cmd.Help()
			return
		}

		tokenName := args[0]

		fmt.Printf("tokenName: %s\n", tokenName)
		generateCmdRun(tokenName)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}

func generateCmdRun(tokenName string) {
	_, tokens, err := Initialize()
	if err != nil {
		return
	}

	token, err := findToken(tokens, tokenName)
	if err != nil {
		fmt.Printf("Error unable to find token: %v\n", err)
		return
	}
	var code string
	var timeLeft int

	for ok := true; ok; ok = timeLeft < MinTimeLeft {
		code, timeLeft = token.GetTotpCode()
		if !ok && verbose {
			fmt.Printf("Got code but time left < %d\n", MinTimeLeft)
		}
	}

	fmt.Printf("code: %v\n", code)
	fmt.Printf("timeLeft: %v\n", timeLeft)
}
