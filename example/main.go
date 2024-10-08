package main

import (
	"log/slog"

	"github.com/nghtf/getmail"

	"os"
)

var log *slog.Logger

func main() {

	log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	var config getmail.TConfig

	config.RCfile = "./getmail.rc"

	gm, err := (&getmail.TGetmail{}).New(log, &config)
	if err != nil {
		log.Error("failed to configure getmail", "error", err)
		return
	}

	err = gm.Fetch()
	if err != nil {
		log.Error("failed to fetch from mailbox", "error", err)
		return
	}

	if err := gm.MailDir.Dispatch(handler); err != nil {
		log.Error("dispatching failed", "error", err)
	}
}

func handler(file string) error {
	log.Info("file recieved", "name", file)
	return nil
}
