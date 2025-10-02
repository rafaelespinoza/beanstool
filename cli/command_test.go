package cli

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
	"testing"

	"github.com/beanstalkd/go-beanstalk"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type CommandSuite struct{}

var _ = Suite(&CommandSuite{})

func (s *CommandSuite) TestCommand_GetStatsForTube(c *C) {
	cmd := &Command{Host: "localhost:11300"}
	err := cmd.Init()
	c.Assert(err, IsNil)

	tube := getRandomTube(cmd.conn)
	tube.Put([]byte(""), 1024, 0, 0)

	stats, err := cmd.GetStatsForTube(tube.Name)
	c.Assert(err, IsNil)
	c.Assert(stats.JobsReady, Equals, 1)
	c.Assert(stats.JobsBuried, Equals, 0)
}

func getRandomTube(conn *beanstalk.Conn) *beanstalk.Tube {
	rb := make([]byte, 32)
	if _, err := rand.Read(rb); err != nil {
		panic(err)
	}

	name := strings.ReplaceAll(base64.URLEncoding.EncodeToString(rb), "=", "0")

	return &beanstalk.Tube{Conn: conn, Name: name}
}
