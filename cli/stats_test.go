package cli

import (
	"bytes"
	"strings"

	"github.com/beanstalkd/go-beanstalk"
	. "gopkg.in/check.v1"
)

type StatsCommandSuite struct {
	c *StatsCommand
	t *beanstalk.Tube
}

var _ = Suite(&StatsCommandSuite{})

func (s *StatsCommandSuite) SetUpTest(c *C) {
	s.c = &StatsCommand{}
	s.c.Host = "localhost:11300"
	err := s.c.Init()
	c.Assert(err, IsNil)

	s.t = getRandomTube(s.c.conn)
	s.c.Tubes = s.t.Name
}

func (s *StatsCommandSuite) TestPutCommand_Put(c *C) {
	s.t.Put([]byte(""), 1024, 0, 0)
	s.t.Put([]byte(""), 1024, 0, 0)
	s.t.Put([]byte(""), 1024, 0, 0)

	var buf bytes.Buffer
	err := s.c.PrintStats(&buf)
	c.Assert(err, IsNil)

	tubeNames, err := s.c.getTubesName()
	c.Assert(err, IsNil)
	c.Assert(len(tubeNames), Equals, 1)
	c.Assert(strings.Contains(buf.String(), tubeNames[0]), Equals, true)
}
