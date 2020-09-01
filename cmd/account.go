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
)

// accountCmd represents the account command
var (
	countrycode, mobile, password string

	accountCmd = &cobra.Command{
		Use:   "account",
		Short: "Authy account info or register device",
		Long: `Register device or show registered account info. 

Can specify country code, mobile number and authy main password.
If not provided, will get from command line stdin`,
		Run: func(cmd *cobra.Command, args []string) {
			registerOrGetDeviceInfo()
		},
	}
)

func init() {
	rootCmd.AddCommand(accountCmd)

	accountCmd.Flags().StringVarP(&countrycode, "countrycode", "c", "", "phone number country code (e.g. 1 for United States), digitals only")
	accountCmd.Flags().StringVarP(&mobile, "mobilenumber", "m", "", "phone number, digitals only")
	accountCmd.Flags().StringVarP(&password, "password", "p", "", "authy main password")
}

func registerOrGetDeviceInfo() {
	_, tokens, err := Initialize()
	if err != nil {
		return
	}

	fmt.Printf("\nLoaded %d auth tokens from authy server\n\n", len(tokens))
	for _, v := range tokens {
		fmt.Printf("Token: %s\n", v.Name)
	}

}
