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
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// package-level constants
const (
	listenPort = "8000"
)

// Thing ...
type Thing struct {
	Name       string
	Activities map[string][]string
}

// Thing.addActions()
func (t Thing) addActions(activities []string) {
	for _, action := range activities {
		t.Activities[action] = append(
			t.Activities[action], time.Now().Format(time.RFC3339))
	}
}

// TimeStats ...
type TimeStats struct {
	Today uint64
	Week  uint64
	Month uint64
	Year  uint64
}

// Stats object ...
type Stats struct {
	Name       string
	Activities map[string]TimeStats
}

// Thing, the JSON version
type jsonThing struct {
	Activities []string `json:"activities"`
}

var memoData = make(map[string]Thing)
var memoStats = make(map[string]Stats)
var mu sync.Mutex

func main() {
	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/", defaultHandler)
	r.HandleFunc("/{thing}", thingPostHandler).Methods("POST")
	r.HandleFunc("/{thing}", thingGetHandler).Methods("GET")

	// Bind to a port and pass our router in
	http.ListenAndServe(":"+listenPort, r)
}

// handler echoes the HTTP request.
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(getAllStats())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// thingPostHandler handles /{thing} GET requests
func thingGetHandler(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["thing"]
	mu.Lock()
	stats, ok := getStats(name)
	mu.Unlock()
	if !ok {
		http.NotFound(w, r)
		return
	}
	js, err := json.Marshal(stats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// thingPostHandler handles /{thing} POST requests
func thingPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "expected application/json", http.StatusUnsupportedMediaType)
		return
	}
	var jsonData jsonThing
	// decode JSON POST body data
	err := json.NewDecoder(r.Body).Decode(&jsonData)
	// grab the name from the route variable, and activities from the POST body
	name, activities := mux.Vars(r)["thing"], jsonData.Activities
	if err != nil {
		log.Print(err)
		return
	}
	mu.Lock()
	thing, _ := getThing(name)
	defer calcStats(thing)
	thing.addActions(activities)
	memoData[name] = thing
	mu.Unlock()
}

// simple data get'er
func getData() map[string]Thing {
	return memoData
}

// make a new Thing struct
func newThing(name string) Thing {
	return Thing{Name: name, Activities: make(map[string][]string)}
}

// get a thing, returns a new thing that the caller can optionally use
func getThing(name string) (Thing, bool) {
	thing, ok := memoData[name]
	if !ok {
		return newThing(name), false
	}
	return thing, true
}

func getAllStats() map[string]Stats {
	return memoStats
}

// make a new Stats struct
func newStats(name string) Stats {
	return Stats{Name: name, Activities: make(map[string]TimeStats)}
}

// simple stats get'er
func getStats(name string) (Stats, bool) {
	stats, ok := memoStats[name]
	if !ok {
		return newStats(name), false
	}
	return stats, true
}

// calculate stats
// TODO: optimize!
func calcStats(thing Thing) {
	stats, _ := getStats(thing.Name)
	for action, times := range thing.Activities {
		var today, week, month, year uint64
		for _, t := range times {
			timestamp, err := time.Parse(time.RFC3339, t)
			if err != nil {
				log.Print(err)
			}
			if isToday(timestamp) {
				today++
			}
			if isThisWeek(timestamp) {
				week++
			}
			if isThisMonth(timestamp) {
				month++
			}
			if isThisYear(timestamp) {
				year++
			}
		}
		timeStats := TimeStats{Today: today, Week: week, Month: month, Year: year}
		stats.Activities[action] = timeStats
	}
	memoStats[thing.Name] = stats
}

// accept a time struct, does it represent a time today?
func isToday(t time.Time) bool {
	tz := t.Location()
	year, month, day := t.Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, tz)
	return t.After(today)
}

// accept a time struct, does it represent a time in the last 7 days?
func isThisWeek(t time.Time) bool {
	oneWeekAgo := time.Now().AddDate(0, 0, -7)
	return t.After(oneWeekAgo)
}

// accept a time struct, does it represent a time in the last month?
func isThisMonth(t time.Time) bool {
	oneMonthAgo := time.Now().AddDate(0, -1, 0)
	return t.After(oneMonthAgo)
}

// accept a time struct, does it represent a time in the last year?
func isThisYear(t time.Time) bool {
	oneYearAgo := time.Now().AddDate(-1, 0, 0)
	return t.After(oneYearAgo)
}
