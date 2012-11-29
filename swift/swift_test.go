package swift_test

import (
	"flag"
	. "launchpad.net/gocheck"
	"launchpad.net/goose/client"
	"launchpad.net/goose/identity"
	"launchpad.net/goose/swift"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

var live = flag.Bool("live", false, "Include live OpenStack (Canonistack) tests")

type SwiftSuite struct {
	swift *swift.Client
}

func (s *SwiftSuite) SetUpSuite(c *C) {
	if !*live {
		c.Skip("-live not provided")
	}

	cred, err := identity.CompleteCredentialsFromEnv()
	c.Assert(err, IsNil)
	client := client.NewOpenStackClient(cred, identity.AuthUserPass)
	err = client.Authenticate()
	c.Assert(err, IsNil)
	c.Logf("client authenticated")
	s.swift = swift.New(client)
}

var suite = Suite(&SwiftSuite{})

func (s *SwiftSuite) TestCreateAndDeleteContainer(c *C) {
	container := "test_container"
	err := s.swift.CreateContainer(container)
	c.Check(err, IsNil)
	err = s.swift.DeleteContainer(container)
	c.Check(err, IsNil)
}

func (s *SwiftSuite) TestObjects(c *C) {

	container := "test_container"
	object := "test_obj"
	data := "...some data..."
	err := s.swift.CreateContainer(container)
	c.Check(err, IsNil)
	err = s.swift.PutObject(container, object, []byte(data))
	c.Check(err, IsNil)
	objdata, err := s.swift.GetObject(container, object)
	c.Check(err, IsNil)
	c.Check(string(objdata), Equals, data)
	err = s.swift.DeleteObject(container, object)
	c.Check(err, IsNil)
	err = s.swift.DeleteContainer(container)
	c.Check(err, IsNil)
}
