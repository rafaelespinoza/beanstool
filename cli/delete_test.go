package cli

import (
	"github.com/beanstalkd/go-beanstalk"
	. "gopkg.in/check.v1"
)

type DeleteCommandSuite struct {
	c *DeleteCommand
	t *beanstalk.Tube
}

var _ = Suite(&DeleteCommandSuite{})

func (s *DeleteCommandSuite) SetUpTest(c *C) {
	s.c = &DeleteCommand{}
	s.c.Host = "localhost:11300"
	err := s.c.Init()
	c.Assert(err, IsNil)

	s.t = getRandomTube(s.c.conn)
	s.c.Tube = s.t.Name
}

func (s *DeleteCommandSuite) TestDeleteReady(c *C) {
	s.t.Put([]byte(""), 1024, 0, 0)
	s.t.Put([]byte(""), 1024, 0, 0)
	s.t.Put([]byte(""), 1024, 0, 0)

	stats, _ := s.c.GetStatsForTube(s.c.Tube)
	c.Assert(stats.JobsReady, Equals, 3)

	s.c.Empty = false
	s.c.State = "ready"

	err := s.c.Delete()
	c.Assert(err, IsNil)

	stats, _ = s.c.GetStatsForTube(s.c.Tube)
	c.Assert(stats.JobsReady, Equals, 2)
}

func (s *DeleteCommandSuite) TestDeleteReadyEmpty(c *C) {
	s.t.Put([]byte(""), 1024, 0, 0)
	s.t.Put([]byte(""), 1024, 0, 0)
	s.t.Put([]byte(""), 1024, 0, 0)

	stats, _ := s.c.GetStatsForTube(s.c.Tube)
	c.Assert(stats.JobsReady, Equals, 3)

	s.c.Empty = true
	s.c.State = "ready"

	err := s.c.Delete()
	c.Assert(err, IsNil)

	stats, _ = s.c.GetStatsForTube(s.c.Tube)
	c.Assert(stats.JobsReady, Equals, 0)
}
