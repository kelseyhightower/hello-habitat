// Copyright 2017 Google Inc. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
//
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	configFile string
)

type Config struct {
	Hostname string `json:"hostname"`
	MaxProcs int    `json:"go_max_procs"`
}

func main() {
	flag.StringVar(&configFile, "config-file", "", "The configuration file")
	flag.Parse()

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("Error getting hostname")
	}

	log.Println("Starting the hello-habitat service...")

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Error loading configuration file: %v", err)
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error loading configuration file: %v", err)
	}

	config.Hostname = hostname

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello Habitat!\n")
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello Habitat Version: 1.0.0\n")
	})

	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		data, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			log.Printf("error formatting JSON response: %v", err)
			http.Error(w, "Error loading the internal configuration", 500)
			return
		}
		w.Write(data)
		fmt.Fprintf(w, "\n")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
