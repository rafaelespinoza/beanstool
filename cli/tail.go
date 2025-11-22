package cli

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/beanstalkd/go-beanstalk"
)

var TooManyErrorsError = errors.New("Too many errors")

type TailCommand struct {
	Tube   string `short:"t" long:"tube" description:"tube to be tailed." required:"true"`
	Action string `short:"a" long:"action" description:"action to perform after reserving the job. (release, bury, delete)" default:"release"`

	Command
}

func (c *TailCommand) Execute(args []string) error {
	if err := c.Init(); err != nil {
		return err
	}

	return c.Tail()
}

func (c *TailCommand) Tail() error {
	ts := beanstalk.NewTubeSet(c.conn, c.Tube)
	lgr := newLogger("tail", c.Tube)
	errors := 0
	for {
		if errors > 100 {
			return TooManyErrorsError
		}

		id, body, err := ts.Reserve(time.Hour * 24)
		if err != nil {
			if err.Error() != "reserve-with-timeout: deadline soon" {
				errors++
				lgr.Warn("Encountered error while reserving, continuing on", slog.String("error", err.Error()))
			}

			continue
		}

		if err := c.PrintJob(id, body); err != nil {
			errors++
			lgr.Warn("Encountered error while printing job, continuing on", slog.String("error", err.Error()), slog.Uint64("job_id", id))
			continue
		}

		if err := c.postPrintAction(id); err != nil {
			return err
		}

		fmt.Println(strings.Repeat("-", 80))
	}
}

func (c *TailCommand) postPrintAction(id uint64) error {
	var err error

	switch c.Action {
	case "release":
		err = c.conn.Release(id, 1024, 0)
	case "bury":
		err = c.conn.Bury(id, 1024)
	case "delete":
		err = c.conn.Delete(id)
	}

	return err
}
