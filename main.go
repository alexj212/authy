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

package main

import (
	"github.com/alexj212/authy/cmd"
)

var (
	// Version version tag
	Version = "v0.13"

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

func main() {
	cmd.Version = Version
	cmd.BuildDate = BuildDate
	cmd.LatestCommit = LatestCommit
	cmd.BuildNumber = BuildNumber
	cmd.BuiltOnIP = BuiltOnIP
	cmd.BuiltOnOs = BuiltOnOs
	cmd.RuntimeVer = RuntimeVer

	cmd.Execute()
}
