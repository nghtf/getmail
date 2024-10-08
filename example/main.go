package main

import (
	"log/slog"

	"github.com/nghtf/getmail"

	"os"
)

func main() {

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	var Getmail getmail.TConfig

	getmail, err := getmail.New(log, &Getmail)
	if err != nil {
		log.Error("getmail.New() failed", "error", err)
	}

	err := getmail.Fetch()
	if err != nil {
		log.Error("getmail.Fetch() failed", "error", err)
	}
	if err := getmail.MailDir.Dispatch(handler); err != nil {
		log.Error("getmail.MailDir.Dispatch() failed", "error", err)

	}

}

func handler(file string) error {
	log.Info("recieved", "filename", file)
	return nil
}
