package cli

import (
	"log/slog"

	"github.com/beanstalkd/go-beanstalk"
)

type KickCommand struct {
	Tube string `short:"t" long:"tube" description:"tube to kick jobs in." required:"true"`
	Num  int    `short:"n" long:"num" description:"number of jobs to kick."`
	Command
}

func (c *KickCommand) Execute(args []string) error {
	if err := c.Init(); err != nil {
		return err
	}

	return c.Kick()
}

func (c *KickCommand) Kick() error {
	if err := c.calcNumIfNeeded(); err != nil {
		return err
	}

	lgr := newLogger("kick", c.Tube)

	if c.Num == 0 {
		lgr.Info("Empty buried queue at tube")
		return nil
	}

	lgr.Info("Trying to kick jobs", slog.Int("num_jobs", c.Num))

	t := beanstalk.NewTube(c.conn, c.Tube)
	kicked, err := t.Kick(c.Num)
	if err != nil {
		return err
	}

	lgr.Info("Actually kicked", slog.Int("num_jobs", kicked))
	return nil
}

func (c *KickCommand) calcNumIfNeeded() error {
	if c.Num == 0 {
		s, err := c.GetStatsForTube(c.Tube)
		if err != nil {
			return err
		}

		c.Num = s.JobsBuried
	}

	return nil
}
