package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	log "unknwon.dev/clog/v2"
)

func main() {
	defer log.Stop()
	err := log.NewConsole()
	if err != nil {
		panic(err)
	}

	// Inspire the world here!
	_ = func() error {
		hostName, err := os.Hostname()
		if err != nil {
			// We use `"github.com/pkg/errors"` here.
			return errors.Wrap(err, "get host name")
		}

		log.Info("Your hostname is %q", hostName)
		return nil
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
}
