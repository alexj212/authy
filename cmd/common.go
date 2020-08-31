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
	"encoding/json"
	"fmt"
	"github.com/alexj212/authy/totp"
	"github.com/alexzorin/authy"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"strings"
	"time"
)

const cacheFileName = ".authycache.json"

var verbose bool

var (
	// BuildDate date string of when build was performed filled in by -X compile flag
	BuildDate string

	// LatestCommit date string of when build was performed filled in by -X compile flag
	LatestCommit string

	// BuildNumber date string of when build was performed filled in by -X compile flag
	BuildNumber string

	// BuiltOnIP date string of when build was performed filled in by -X compile flag
	BuiltOnIP string

	// BuiltOnOs date string of when build was performed filled in by -X compile flag
	BuiltOnOs string

	// RuntimeVer date string of when build was performed filled in by -X compile flag
	RuntimeVer string
)

// Token save in cache
type Token struct {
	Name         string `json:"name"`
	OriginalName string `json:"original_name"`
	Digital      int    `json:"digital"`
	Secret       string `json:"secret"`
	Period       int    `json:"period"`
}

// Tokens type for results of search etc
type Tokens []*Token

// String return token name of results index
func (ts Tokens) String(i int) string {
	if len(ts[i].Name) > len(ts[i].OriginalName) {
		return ts[i].Name
	}

	return ts[i].OriginalName
}

// Len - number of Token results
func (ts Tokens) Len() int { return len(ts) }

func loadCachedTokens() (tks []*Token, err error) {
	fpath, err := ConfigPath(cacheFileName)
	if err != nil {
		return
	}

	f, err := os.Open(fpath)
	if err != nil {
		return
	}

	defer f.Close()
	err = json.NewDecoder(f).Decode(tks)
	if verbose {
		fmt.Printf("Loaded cached providers from %v\n", fpath)
	}
	return
}

func saveTokens(tks []*Token) (err error) {
	regrPath, err := ConfigPath(cacheFileName)
	if err != nil {
		return
	}

	f, err := os.OpenFile(regrPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return
	}

	defer f.Close()
	err = json.NewEncoder(f).Encode(tks)

	if err != nil {
		fmt.Printf("Error saving tokens: %v\n", err)
		return
	}
	if verbose {
		fmt.Printf("Saved tokens to file: %v\n", regrPath)
	}
	return
}

func getTokensFromAuthyServer(devInfo *DeviceRegistration) (tks []*Token, err error) {
	client, err := authy.NewClient()
	if err != nil {
		log.Fatalf("Create authy API client failed %+v", err)
	}

	apps, err := client.QueryAuthenticatorApps(nil, devInfo.UserID, devInfo.DeviceID, devInfo.Seed)
	if err != nil {
		log.Fatalf("Fetch authenticator apps failed %+v", err)
	}

	if !apps.Success {
		log.Fatalf("Fetch authenticator apps failed %+v", apps)
	}

	tokens, err := client.QueryAuthenticatorTokens(nil, devInfo.UserID, devInfo.DeviceID, devInfo.Seed)
	if err != nil {
		log.Fatalf("Fetch authenticator tokens failed %+v", err)
	}

	if !tokens.Success {
		log.Fatalf("Fetch authenticator tokens failed %+v", tokens)
	}

	if len(devInfo.MainPassword) == 0 {
		fmt.Print("\nPlease input Authy main password: ")
		pp, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			log.Fatalf("Get password failed %+v", err)
		}

		devInfo.MainPassword = strings.TrimSpace(string(pp))
		SaveDeviceInfo(*devInfo)
	}

	tks = []*Token{}
	for _, v := range tokens.AuthenticatorTokens {
		if verbose {
			fmt.Printf("AuthenticatorTokens: %v\n", v.Name)
		}

		secret, err := v.Decrypt(devInfo.MainPassword)
		if err != nil {
			//fmt.Printf("v.Name: %v error: %v\n", v.Name, err)
			log.Fatalf("Decrypt token failed %+v", err)
		}

		tks = append(tks, &Token{
			Name:         v.Name,
			OriginalName: v.OriginalName,
			Digital:      v.Digits,
			Secret:       secret,
		})
	}

	for _, v := range apps.AuthenticatorApps {
		secret, err := v.Token()
		if err != nil {
			log.Fatal("Get secret from app failed", err)
		}
		if verbose {
			fmt.Printf("AuthenticatorApps: %v\n", v.Name)
		}

		tks = append(tks, &Token{
			Name:    v.Name,
			Digital: v.Digits,
			Secret:  secret,
			Period:  10,
		})
	}

	saveTokens(tks)
	return
}

func findToken(tokenName string) (*Token, error) {
	devInfo, err := LoadExistingDeviceInfo()
	if err != nil {
		if os.IsNotExist(err) {
			devInfo, err = newRegistrationDevice()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	tokens, err := loadCachedTokens()
	if err != nil {
		tokens, err = getTokensFromAuthyServer(&devInfo)
		if err != nil {
			err = errors.Wrap(err, "unable to getTokensFromAuthyServer")
			return nil, err
		}
	}

	var tk *Token
	for _, v := range tokens {
		if tokenName == v.Name {
			tk = v
			break
		}
	}

	if tk == nil {
		err = errors.Errorf("unable to find token: %s", tokenName)
		return nil, err
	}
	return tk, nil
}

// GetTotpCode return code and # of seconds left in current code
func (tk *Token) GetTotpCode() (string, int) {
	codes := totp.GetTotpCode(tk.Secret, tk.Digital)
	challenge := totp.GetChallenge()

	secsLeft := 30 - int(time.Now().Unix()-challenge*30)
	code := codes[1]
	return code, secsLeft
}
