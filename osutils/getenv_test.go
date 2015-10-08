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

package osutils

import (
	"os"
	"testing"
)

func TestGetenv(t *testing.T) {
	if err := os.Setenv("testEnvVar", "venkman"); err != nil {
		t.Error("Something is wrong with the built-in os package!")
	}

	if testEnvVar := Getenv("testEnvVar", "default"); testEnvVar != "venkman" {
		t.Error("Unable to get a known environment variable!")
	}

	os.Setenv("novelEnvVar", "") // make sure this is "unset" before testing for default
	if novelEnvVar := Getenv("novelEnvVar", "spengler"); novelEnvVar != "spengler" {
		t.Error("Didn't get the default passed-in string as expected!")
	}

}
