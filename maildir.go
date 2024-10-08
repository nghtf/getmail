package getmail

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
)

const DIR_NEW = "new"
const DIR_CUR = "cur"
const DIR_TMP = "tmp"
const DIR_SYS = ".getmail"

type TMailDir struct {
	log     *slog.Logger
	path    string
	ready   bool
	folders map[string]string
}

// Initialize maildir and create all subfolders if necessary
func (md *TMailDir) New(log *slog.Logger, path string) (*TMailDir, error) {

	md.log = log

	if path == "" {
		return nil, errors.New("maildir path not specified")
	}
	md.path, _ = filepath.Abs(path)

	dirs := []string{DIR_NEW, DIR_CUR, DIR_TMP, DIR_SYS}
	md.folders = make(map[string]string)

	for _, dir := range dirs {
		md.folders[dir] = filepath.Join(md.path, dir)
		err := os.MkdirAll(md.folders[dir], os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	md.ready = true
	return md, nil
}

type Handler_MaildirDispatcher func(file string) error

// Sends all files from Maildir/New folder to handler. Move each file to Maildir/Cur if handler succeeds.
func (md *TMailDir) Dispatch(handler Handler_MaildirDispatcher) error {

	if !md.ready {
		return errors.New("maildir is not initialized properly")
	}

	items, err := os.ReadDir(md.folders[DIR_NEW])
	if err != nil {
		return err
	}

	log := md.log
	log.Info("found " + strconv.Itoa(len(items)) + " new message(s)")

	for _, item := range items {
		if !item.IsDir() {
			file := filepath.Join(md.folders[DIR_NEW], item.Name())
			if err := handler(file); err != nil {
				return err
			}
			if err = md.Move(file, md.folders[DIR_CUR]); err != nil {
				return err
			}
		}
	}

	return nil
}

func (md *TMailDir) Move(file string, dir string) error {
	err := os.Rename(file, filepath.Join(dir, filepath.Base(file)))
	if err != nil {
		return err
	}
	return nil
}
