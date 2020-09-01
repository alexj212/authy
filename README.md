## Authy
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0) [![GoDoc](https://godoc.org/github.com/alexj212/authy?status.png)](http://godoc.org/github.com/alexj212/authy)  [![travis](https://travis-ci.org/alexj212/authy.svg?branch=master)](https://travis-ci.org/alexj212/authy) [![Go Report Card](https://goreportcard.com/badge/github.com/alexj212/authy)](https://goreportcard.com/report/github.com/alexj212/authy)


Authy Commandline Tool

## Description
`authy` will allow the user to login to authy and cache totp tokens for command line generation. 
Sub commands exists that will allow for generation of a token, listing available tokens, executing a script with generated token. 

#### Installation

Pre-built binaries are available from the [releases page](https://github.com/alexj212/authy/releases).

Alternatively, it can be compiled from source, which requires [Go 1.14 or newer](https://golang.org/doc/install):

```
go get github.com/alexj212/authy
```

#### To use it
1. Move downloaded binary to your local `$PATH`
2. Run `authy account`. The command will prompt you for your phone number country code (e.g. 1 for United States) and your phone number. This is the number that you used to register your Authy account originally. 
3. If the program identifies an existing Authy account, it will send a device registration request using the push method. This will send a push notification to your existing Authy apps (be it on Android, iOS, Desktop or Chrome), and you will need to respond that from your other app(s).
4. If the device registration is successful, the program will save its authentication credential (a random value) to `$HOME/.authy.json` for further uses. It will prompt for the authy account master password. It will attempt to decrypt totp tokens with the master password. If the wrong password is entered it will fail decruption.   
5. Run `authy list` to list available tokens. 

#### Commands
##### authy account configure/display account information 
```bash
$ ./bin/authy account

What is your phone number's country code? (digits only, e.g. 1): 1

What is your phone number? (digits only): XXXXXXXXXX
2020/08/31 16:55:02 Checking device registration status (4m59s until we give up)
2020/08/31 16:55:07 Checking device registration status (4m54s until we give up)
2020/08/31 16:55:12 Checking device registration status (4m49s until we give up)
2020/08/31 16:55:13 Register device success!!!
2020/08/31 16:55:13 Your device info: {UserID:XXXXXXX DeviceID:XXXXXXXXX Seed:XXXXXXXXXXXXXXXXXXXXXXXX APIKey:XXXXXXXXXXXXXXXXXXXXXXXX MainPassword:}

Please input Authy main password:
 
Loaded 2 auth tokens from authy server
Token: alexj@backpocket.com
Token: Twilio

```

##### authy list
list available totp accounts
```bash
    `authy list [REGEX]` list accounts cached from ~/.authycache.json

$ ./bin/authy list
Token: alexj@backpocket.com
Token: Twilio

```

##### authy refresh 
reload available totp accounts from authy
```bash
$ ./bin/authy --verbose refresh
Loaded device info from file: .authy.json
AuthenticatorTokens: alexj@backpocket.com
AuthenticatorApps: Twilio
Saved tokens to file: /home/alexj/.authycache.json
```


##### authy delpwd
 remove cached master passwords
```bash
./bin/authy --verbose delpwd
Loaded device info from file: .authy.json
Save device info to file: /home/alexj/.authy.json
2020/08/31 16:11:10 Backup password delete successfully!

```


##### authy exec
 generate token and invoke script with totp token
```bash

$ ./bin/authy -v exec alexj@backpocket.com "/home/alexj/bin/vpnup.sh [AUTHCODE]" 
tokenName: alexj@backpocket.com
script: /home/alexj/bin/vpnup.sh [AUTHCODE]
replacement: [AUTHCODE]
dryRun: false
Loaded device info from file: .authy.json
Loaded cached providers from /home/alexj/.authycache.json
AuthenticatorTokens: alexj@backpocket.com
Saved tokens to file: /home/alexj/.authycache.json

Script executed 
"/home/alexj/bin/vpnup.sh XXXX
```



##### authy help
 display help
```bash
$ ./bin/authy help
Authy command line tool

Usage:
  authy [command]

Available Commands:
  account     Authy account info or register device
  delpwd      Delete saved backup password
  exec        exec a program/script and pass otp token
  generate    generate a otp token
  help        Help about any command
  info        Display info on authy cmd
  list        list search your otp tokens(case-insensitive)
  refresh     Refresh token cache
  wipe        remove ~/.authy*.json cache files

Flags:
  -h, --help      help for authy
  -v, --verbose   verbose output

Use "authy [command] --help" for more information about a command.

```
     
#### Files
    ~/.authy.json    
    ~/.authycache.json
    
   
#### Building
```bash
$ make authy


$ make help
build_info                     Build the container
help                           This help.
all                            build example and run tests
binaries                       build binaries in bin dir
authy                          build example binary in bin dir
clean                          clean all binaries in bin dir
clean_binary                   clean binary in bin dir
clean_authy                    clean dumper
test                           run tests
fmt                            run fmt on project
doc                            launch godoc on port 6060
deps                           display deps for project
lint                           run lint on the project
staticcheck                    run staticcheck on the project
vet                            run go vet on the project
tools                          install dependent tools for code analysis
gocyclo                        run gocyclo on the project
check                          run code checks on the project

```   
    
#### Inspiration
    https://github.com/momaek/authy
    
#### Author
Alex jeannopoulos, alexj@backpocket.com    

#### LICENSE
[Apache 2.0](./LICENSE)