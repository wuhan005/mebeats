// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	log "unknwon.dev/clog/v2"
)

func main() {
	defer log.Stop()
	err := log.NewConsole()
	if err != nil {
		panic(err)
	}

	deviceName := flag.String("device-name", "", "Mi Band device name.")
	authKey := flag.String("auth-key", "", "Mi Band auth key.")
	flag.Parse()

	log.Trace("Try to connect %q...", deviceName)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
}
