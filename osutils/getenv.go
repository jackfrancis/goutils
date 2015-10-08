/*
Copyright 2015 Jack Francis

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package osutils contains thin wrappers around built-in os packages
package osutils

import (
	"os"
)

// Getenv passes along the 1st arg to os.Getenv and returns its value,
// or a passed-in default (2nd arg) if os.Getenv returns an empty string
func Getenv(envVar string, defaultVal string) string {
	if val := os.Getenv(envVar); len(val) == 0 {
		return defaultVal
	} else {
		return val
	}
}
