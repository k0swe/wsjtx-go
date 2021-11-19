package integration

import (
	"net"
	"testing"

	"github.com/k0swe/wsjtx-go/v3"
	"github.com/stretchr/testify/suite"
)

type integrationTestSuite struct {
	suite.Suite
	stub    wsjtx.Server
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
	s.stub, err = wsjtx.MakeServerGiven(net.ParseIP("127.0.0.1"), 0)
	s.Require().NoError(err)
	go s.stub.ListenToWsjtx(s.msgChan, s.errChan)

}

func (s *integrationTestSuite) SetupTest() {
	var err error
	s.fake, err = NewFake(s.stub.LocalAddr())
	s.Require().NoError(err)
}

func (s *integrationTestSuite) TearDownTest() {
	s.fake = &WsjtxFake{}
}
