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
	"regexp"

	//"time"

	"github.com/spf13/cobra"
)

var (
	// MinTimeLeft minimum number of seconds left when returning a code, api will block till greater than value
	MinTimeLeft = 10
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [TOKEN_PATTERN]...",
	Short: "list search your otp tokens(case-insensitive)",
	Long: `list search your otp tokens(case-insensitive)

First time(or after clean cache) , need your authy main password`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) > 0 {
			for _, pattern := range args {
				listSearch(pattern)
			}

		} else {
			listSearch(".*")
		}

	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func listSearch(patternStr string) {
	pattern, err := regexp.Compile(patternStr)
	if err != nil {
		fmt.Printf("invalid pattern: %s error: %v\n", patternStr, err)
		return
	}

	_, tokens, err := Initialize()
	if err != nil {
		return
	}

	found := false
	for _, v := range tokens {
		if pattern.MatchString(v.Name) {
			fmt.Printf("Token: %s\n", v.Name)
			found = true
		}
	}

	if !found {
		fmt.Printf("unable to match pattern: %s to any tokens\n", pattern)
	}

}
