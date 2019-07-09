// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package main

import (
	"os"

	"github.com/spf13/pflag"
)

var remoteUsername *string = pflag.StringP("username", "u", "", "username for auth")

var remotePassword *string = pflag.StringP("password", "p", "", "password for remote auth")

var remoteURL = os.Args[len(os.Args)-1]

func main() {
	pflag.Parse()
	err := createWindow()
	if err != nil {
		panic(err.Error())
	}

}
