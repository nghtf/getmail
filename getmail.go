package getmail

import (
	"log/slog"
	"os/exec"
	"path/filepath"

	"github.com/nghtf/rcscan"
)

type TConfig struct {
	// REQUIRED. Mailbox configuration file
	RCfile string `yaml:"rcfile" env-required:"true"`
	// OPTIONAL. Folder to store getmail runtime data
	GetmailDir string `yaml:"getmaildir"`
	// OPTIONAL. Folder with the mail recieved. Trailing "/" is required.
	MaildirPath string `yaml:"maildir"`
}

type TGetmail struct {
	log     *slog.Logger
	config  *TConfig
	MailDir *TMailDir
}

// Setup getmail wrapper
func (gm *TGetmail) New(log *slog.Logger, config *TConfig) (*TGetmail, error) {

	var err error

	gm.log = log
	gm.config = config

	// Check if maildir path configured
	if gm.config.MaildirPath == "" {
		log.Info("maildir not specified, pulling from file", "file", gm.config.RCfile)
		rc, err := rcscan.New(gm.config.RCfile)
		if err != nil {
			return nil, err
		} else {
			gm.config.MaildirPath, err = rc.Get("destination", "path")
			if err != nil {
				return nil, err
			}
		}
	}

	// Instantiate maildir object and create all folders if it needs (including default system "maildir/.getmail" directory)
	gm.MailDir, err = (&TMailDir{}).New(log, gm.config.MaildirPath)
	if err != nil {
		return nil, err
	}

	// Make sure getmailDir is there or map it to default path "maildir/.getmail", that created by maildir.New()
	if gm.config.GetmailDir == "" {
		gm.config.GetmailDir = filepath.Join(gm.config.MaildirPath, DIR_SYS)
		log.Info("getmailDir not specified, using default location", "path", gm.config.GetmailDir)
	}

	return gm, nil
}

// Fetches mail from mailbox
func (gm *TGetmail) Fetch() error {

	var err error

	// GetMail requires that to be of an absolute path (otherwise it tries to merge it with GetmailDir)
	gm.config.RCfile, err = filepath.Abs(gm.config.RCfile)
	if err != nil {
		return err
	}

	// External call

	c := exec.Command("getmail", "--rcfile="+gm.config.RCfile, "--getmaildir="+gm.config.GetmailDir)
	out, err := c.CombinedOutput()
	gm.log.Debug("getmail()", "stdout", string(out))
	if err != nil {
		return err
	}
	return nil
}
