// Copyright 2019 Northern.tech AS
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path"
	"path/filepath"
)

func CheckErr(e error) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "FAILURE: %s\n", e.Error())
		os.Exit(1)
	}
}

func migrateAuthToken(oldtoken string, token string) {
	// if needed, migrate token from old to new location
	if _, err := os.Stat(token); !os.IsNotExist(err) {
		// new token exists, no migration
		return
	}

	if _, err := os.Stat(oldtoken); err != nil {
		// old token doesn't exist, no migration
		return
	}

	// Attempt migration, ignore errors (but log them?)
	if err := os.MkdirAll(filepath.Dir(token), 0700); err == nil {
		// log that token was moved?
		err = os.Rename(oldtoken, token)
	}

	// Cleanup old token directory if empty
	os.Remove(filepath.Dir(oldtoken)) // err on non-empty, ignore.
}

func getDefaultAuthTokenPath() (string, error) {
	cachedir := ""
	userhomedir := ""

	if homeenv := os.Getenv("HOME"); homeenv != "" {
		userhomedir = homeenv
	} else if user, err := user.Current(); err == nil {
		userhomedir = user.HomeDir
	} else {
		return "", errors.New("Not able to determine users cache dir. Is the '$HOME', and '$USER' variable a part of the running program environment?")
	}

	if cachehomeenv := os.Getenv("XDG_CACHE_HOME"); cachehomeenv != "" {
		cachedir = cachehomeenv
	} else {
		cachedir = path.Join(userhomedir, ".cache")
	}

	oldtoken := filepath.Join(userhomedir, ".mender", "authtoken")
	token := filepath.Join(cachedir, "mender", "authtoken")

	migrateAuthToken(oldtoken, token)

	return token, nil
}
