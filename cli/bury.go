package cli

import (
	"log/slog"
	"strconv"
	"time"

	"github.com/beanstalkd/go-beanstalk"
)

type BuryCommand struct {
	Tube string `short:"t" long:"tube" description:"tube to bury jobs in." required:"true"`
	Num  int    `short:"n" long:"num" description:"number of jobs to bury."`
	Command
}

func (c *BuryCommand) Execute(args []string) error {
	if err := c.Init(); err != nil {
		return err
	}

	return c.Bury()
}

func (c *BuryCommand) Bury() error {
	if err := c.calcNum(); err != nil {
		return err
	}
	lgr := newLogger("bury", c.Tube)

	if c.Num == 0 {
		lgr.Info("Empty ready queue at tube")
		return nil
	}

	lgr.Info("Trying to bury jobs", slog.Int("num_jobs", c.Num))

	count := 0
	ts := beanstalk.NewTubeSet(c.conn, c.Tube)
	for count < c.Num {
		id, _, err := ts.Reserve(time.Second)
		if err != nil {
			return err
		}

		s, err := c.conn.StatsJob(id)
		if err != nil {
			return err
		}

		pri, err := strconv.ParseUint(s["pri"], 10, 32)
		if err != nil {
			return err
		}

		if err := c.conn.Bury(id, uint32(pri)); err != nil {
			return err
		}

		count++
	}

	lgr.Info("Actually buried", slog.Int("num_jobs", count))
	return nil
}

func (c *BuryCommand) calcNum() error {
	s, err := c.GetStatsForTube(c.Tube)
	if err != nil {
		return err
	}

	if c.Num == 0 || c.Num > s.JobsReady {
		c.Num = s.JobsReady
	}

	return nil
}
