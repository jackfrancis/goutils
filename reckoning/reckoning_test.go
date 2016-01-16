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

package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func newServer() *httptest.Server {
	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/", defaultHandler)
	r.HandleFunc("/{thing}", thingPostHandler).Methods("POST")
	r.HandleFunc("/{thing}", thingGetHandler).Methods("GET")
	return httptest.NewServer(r)
}

func TestGetSlash(t *testing.T) {
	server := newServer()
	defer server.Close()
	resp, err := httpGet(server, "")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}

	//actual, err := ioutil.ReadAll(resp.Body)
}

func TestPost(t *testing.T) {
	packageName := "mypackage"
	jsonData := `{"activities": ["install"]}`
	server := newServer()
	defer server.Close()
	resp, err := httpPost(server.URL+"/"+packageName, jsonData)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
	resp, err = httpGet(server, "")
	defer resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
	stats := parseJSON(resp)
	if stats[packageName].Name != packageName {
		t.Error("unexpected package name from JSON response")
	}
	if stats[packageName].Activities["install"].Today != 1 {
		t.Error("unexpected value for Today in JSON response")
	}
	if stats[packageName].Activities["install"].Week != 1 {
		t.Error("unexpected value for Week in JSON response")
	}
	if stats[packageName].Activities["install"].Month != 1 {
		t.Error("unexpected value for Month in JSON response")
	}
	if stats[packageName].Activities["install"].Year != 1 {
		t.Error("unexpected value for Year in JSON response")
	}
}

func parseJSON(r *http.Response) map[string]Stats {
	rawJSON, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
	}
	var rawJSONMap map[string]*json.RawMessage
	err = json.Unmarshal(rawJSON, &rawJSONMap)
	if err != nil {
		log.Print(err)
	}
	stats := make(map[string]Stats)
	for k := range rawJSONMap {
		var statsObj Stats
		err = json.Unmarshal(*rawJSONMap[k], &statsObj)
		stats[k] = statsObj
	}
	return stats
}

func httpGet(s *httptest.Server, route string) (*http.Response, error) {
	return http.Get(s.URL + route)
}

func httpPost(url string, json string) (*http.Response, error) {
	jsonStr := []byte(json)
	return http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
}
