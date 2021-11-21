package integration

import (
	"net"
	"strconv"
	"testing"

	"github.com/k0swe/wsjtx-go/v4"
	"github.com/stretchr/testify/suite"
)

type integrationTestSuite struct {
	suite.Suite
	server  wsjtx.Server
	msgChan chan interface{}
	errChan chan error
	fake    *WsjtxFake
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, &integrationTestSuite{})
}

func (s *integrationTestSuite) SetupSuite() {
	var err error
	s.msgChan = make(chan interface{}, 5)
	s.errChan = make(chan error, 5)
	s.server, err = wsjtx.MakeServerGiven(net.ParseIP("127.0.0.1"), 0)
	s.Require().NoError(err)
	go s.server.ListenToWsjtx(s.msgChan, s.errChan)
	s.T().Log("suite started server listening")
}

func (s *integrationTestSuite) SetupTest() {
	_, portStr, _ := net.SplitHostPort(s.server.LocalAddr().String())
	port, _ := strconv.Atoi(portStr)
	var err error
	s.fake, err = NewFake(&net.UDPAddr{Port: port}, s.T())
	s.Require().NoError(err)
	s.T().Log("suite reports fake is connected")
}

func (s *integrationTestSuite) TearDownTest() {
	s.fake.Stop()
}
