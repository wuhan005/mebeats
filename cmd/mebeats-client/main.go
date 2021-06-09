// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/hex"
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"

	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/mebeats/miband"
	"github.com/wuhan005/mebeats/report"
)

func main() {
	defer log.Stop()
	err := log.NewConsole()
	if err != nil {
		panic(err)
	}

	addr := flag.String("addr", "", "Mi Band device address.")
	key := flag.String("auth-key", "", "Mi Band auth key.")
	serverAddr := flag.String("server-addr", "", "The server address of mebeats.")
	flag.Parse()

	deviceAddr := strings.ToLower(*addr)
	deviceAuthKey, err := hex.DecodeString(strings.ReplaceAll(*key, "-", ""))
	if err != nil {
		log.Fatal("Failed to decode auth key: %v", err)
	}

	log.Trace("Try to connect %q...", deviceAddr)

	band, err := miband.NewMiBand(deviceAddr, string(deviceAuthKey))
	if err != nil {
		log.Fatal("Failed to new Mi Band: %v", err)
	}
	err = band.Initialize()
	if err != nil {
		log.Fatal("Failed to initialize Mi Band: %v", err)
	}

	err = band.GetHeartRateOneTime()
	if err != nil {
		log.Fatal("Failed to init heart rate: %v", err)
	}

	if *serverAddr != "" {
		// Report to server.
		go func() {
			ch := band.Subscribe()
			for {
				select {
				case <-ch:
					err := report.ToServer(*serverAddr,
						report.Options{
							HeartRate: band.GetCurrentHeartRate(),
						},
					)
					if err != nil {
						log.Error("Failed to report to server: %v", err)
					}
				}
			}
		}()
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
}
