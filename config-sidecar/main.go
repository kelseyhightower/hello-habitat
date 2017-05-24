// Copyright 2017 Google Inc. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
//
// You may obtain a copy of the License at
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
)

var (
	configFile    string
	configVersion int
	serviceGroup  string
	watcher       *fsnotify.Watcher
)

func main() {
	flag.StringVar(&configFile, "user-configuration-file", "/etc/habitat/user.toml", "Path to service configuration")
	flag.StringVar(&serviceGroup, "service-group", "", "The Habitat service group")
	flag.Parse()

	configVersion = 1

	log.Println("Starting the Habitat config sidecar...")

	if err := loadConfig(); err != nil {
		log.Fatal(err)
	}

	if err := newWatcher(); err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			if watcher == nil {
				if err := resetWatcher(); err != nil {
					log.Println(err)
				}
				time.Sleep(5 * time.Second)
				continue
			}

			select {
			case <-watcher.Events:
				if err := loadConfig(); err != nil {
					log.Println(err)
				}
				if err := resetWatcher(); err != nil {
					log.Println(err)
				}
			case err := <-watcher.Errors:
				log.Println(err)
			}
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Printf("Shutdown signal received, exiting...")
}

func newWatcher() error {
	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	return watcher.Add(configFile)
}

func resetWatcher() error {
	err := watcher.Close()
	if err != nil {
		return err
	}
	return newWatcher()
}

func loadConfig() error {
	out, err := exec.Command("hab", "config", "apply", serviceGroup,
		strconv.Itoa(configVersion), configFile).CombinedOutput()
	if err != nil {
		log.Println(string(out))
		return err
	}

	log.Println(string(out))
	configVersion += 1

	return nil
}
