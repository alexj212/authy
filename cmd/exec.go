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
	"os"
	"os/exec"
	"strings"
	//"time"

	"github.com/spf13/cobra"
)

// ./bin/authy exec PalVPN "/home/alexj/bin/vpnup.sh [AUTHCODE]"

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec [TokenName] [SCRIPT]",
	Short: "exec a program/script and pass otp token",
	Long: `exec a program/script and pass otp token

First time(or after clean cache) , need your authy main password`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 2 {
			cmd.Help()
			return
		}

		tokenName := args[0]
		script := args[1]
		replacement, err := cmd.Flags().GetString("replacement")
		if err != nil {
			cmd.Help()
			return
		}

		dryRun, err := cmd.Flags().GetBool("dry-run")
		if err != nil {
			cmd.Help()
			return
		}

		fmt.Printf("tokenName: %s\n", tokenName)
		fmt.Printf("script: %s\n", script)
		fmt.Printf("replacement: %v\n", replacement)
		fmt.Printf("dryRun: %v\n", dryRun)
		execCmdRun(tokenName, script, replacement, dryRun)
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
	execCmd.Flags().StringP("replacement", "r", "[AUTHCODE]", "replacement value to substitute auth code in script")
	execCmd.Flags().BoolP("dry-run", "n", false, "dry run")
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func execCmdRun(tokenName, script, replacementToken string, dryRun bool) {
	_, tokens, err := Initialize()
	if err != nil {
		return
	}

	token, err := findToken(tokens, tokenName)
	if err != nil {
		fmt.Printf("Error unable to find token: %v\n", err)
		return
	}

	args := strings.Split(script, " ")

	if !fileExists(args[0]) {
		fmt.Print("exec script: %s does not exist\n")
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

	for i, val := range args {
		args[i] = strings.Replace(val, replacementToken, code, -1)
	}

	if dryRun || verbose {
		fmt.Printf("code: %v\n", code)
		fmt.Printf("timeLeft: %v\n", timeLeft)
		fmt.Printf("orig script: %v\n", script)
		fmt.Printf("script: [%v]\n", strings.Join(args, " "))
	}

	if dryRun {
		fmt.Printf("dry run exiting\n")
		return
	}
	// construct `go version` command
	cmdGoVer := &exec.Cmd{
		Path:   args[0],
		Args:   args,
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}

	if err := cmdGoVer.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}
