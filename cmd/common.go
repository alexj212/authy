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
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/alexj212/authy/totp"
	"github.com/alexzorin/authy"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	configFileName = ".authy.json"
	cacheFileName  = ".authycache.json"
)

var verbose bool

var (
	// Version version tag
	Version = "v0.1"
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

// DeviceRegistration authy account details
type DeviceRegistration struct {
	UserID       uint64 `json:"user_id,omitempty"`
	DeviceID     uint64 `json:"device_id,omitempty"`
	Seed         string `json:"seed,omitempty"`
	APIKey       string `json:"api_key,omitempty"`
	MainPassword string `json:"main_password,omitempty"`
}

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
		secret, err := v.Decrypt(devInfo.MainPassword)
		if err != nil {
			devInfo.MainPassword = ""
			SaveDeviceInfo(*devInfo)
			log.Fatalf("Decrypt token failed %+v", err)
		}

		if verbose {
			fmt.Printf("AuthenticatorTokens: %v\n", v.Name)
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

// SaveDeviceInfo ..
func SaveDeviceInfo(devInfo DeviceRegistration) (err error) {
	regrPath, err := ConfigPath(configFileName)
	if err != nil {
		return
	}

	f, err := os.OpenFile(regrPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return
	}

	defer f.Close()
	err = json.NewEncoder(f).Encode(devInfo)
	if verbose {
		fmt.Printf("Save device info to file: %s\n", regrPath)
	}
	return
}

// LoadExistingDeviceInfo ,,,
func LoadExistingDeviceInfo() (devInfo DeviceRegistration, err error) {
	devPath, err := ConfigPath(configFileName)
	if err != nil {
		log.Println("Get device info file path failed", err)
		os.Exit(1)
	}

	f, err := os.Open(devPath)
	if err != nil {
		return
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&devInfo)

	if err == nil {
		if verbose {
			fmt.Printf("Loaded device info from file: %s\n", configFileName)
		}
	}
	return
}

// ConfigPath get config file path
func ConfigPath(fname string) (string, error) {
	devPath, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	return filepath.Join(devPath, fname), nil
}

func newRegistrationDevice() (devInfo DeviceRegistration, err error) {
	var (
		sc      = bufio.NewScanner(os.Stdin)
		phoneCC int
	)

	if len(countrycode) == 0 {
		fmt.Print("\nWhat is your phone number's country code? (digits only, e.g. 1): ")
		if !sc.Scan() {
			err = errors.New("Please provide a phone country code, e.g. 1")
			log.Println(err)
			return
		}

		countrycode = sc.Text()
	}

	phoneCC, err = strconv.Atoi(strings.TrimSpace(countrycode))
	if err != nil {
		log.Println("Invalid country code. Parse country code failed", err)
		return
	}

	if len(mobile) == 0 {
		fmt.Print("\nWhat is your phone number? (digits only): ")
		if !sc.Scan() {
			err = errors.New("Please provide a phone number, e.g. 1232211")
			log.Println(err)
			return
		}

		mobile = sc.Text()
	}

	mobile = strings.TrimSpace(mobile)

	client, err := authy.NewClient()
	if err != nil {
		log.Println("New authy client failed", err)
		return
	}

	userStatus, err := client.QueryUser(nil, phoneCC, mobile)
	if err != nil {
		log.Println("Query user failed", err)
		return
	}

	if !userStatus.IsActiveUser() {
		err = errors.New("There doesn't seem to be an Authy account attached to that phone number")
		log.Println(err)
		return
	}

	// Begin a device registration using Authy app push notification
	regStart, err := client.RequestDeviceRegistration(nil, userStatus.AuthyID, authy.ViaMethodPush)
	if err != nil {
		log.Println("Start register device failed", err)
		return
	}

	if !regStart.Success {
		err = fmt.Errorf("Authy did not accept the device registration request: %+v", regStart)
		log.Println(err)
		return
	}

	var regPIN string
	timeout := time.Now().Add(5 * time.Minute)
	for {
		if timeout.Before(time.Now()) {
			err = errors.New("Gave up waiting for user to respond to Authy device registration request")
			log.Println(err)
			return
		}

		log.Printf("Checking device registration status (%s until we give up)", time.Until(timeout).Truncate(time.Second))

		regStatus, err1 := client.CheckDeviceRegistration(nil, userStatus.AuthyID, regStart.RequestID)
		if err1 != nil {
			err = err1
			log.Println(err)
			return
		}
		if regStatus.Status == "accepted" {
			regPIN = regStatus.PIN
			break
		} else if regStatus.Status != "pending" {
			err = fmt.Errorf("Invalid status while waiting for device registration: %s", regStatus.Status)
			log.Println(err)
			return
		}

		time.Sleep(5 * time.Second)
	}

	regComplete, err := client.CompleteDeviceRegistration(nil, userStatus.AuthyID, regPIN)
	if err != nil {
		log.Println(err)
		return
	}

	if regComplete.Device.SecretSeed == "" {
		err = errors.New("Something went wrong completing the device registration")
		log.Println(err)
		return
	}

	devInfo = DeviceRegistration{
		UserID:   regComplete.AuthyID,
		DeviceID: regComplete.Device.ID,
		Seed:     regComplete.Device.SecretSeed,
		APIKey:   regComplete.Device.APIKey,
	}

	if verbose {
		fmt.Printf("APIKey       : %v\n", devInfo.APIKey)
		fmt.Printf("DeviceID     : %v\n", devInfo.DeviceID)
		fmt.Printf("MainPassword : %v\n", devInfo.MainPassword)
		fmt.Printf("UserID       : %v\n", devInfo.UserID)
		fmt.Printf("Seed         : %v\n", devInfo.Seed)
	}

	err = SaveDeviceInfo(devInfo)
	if err != nil {
		log.Println("Save device info failed", err)
	}

	return
}
