package integration

import (
	"encoding/hex"
	"net"
	"testing"
	"time"

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

	s.fake, err = NewFake(s.stub.LocalAddr())
	s.Require().NoError(err)
}

func (s *integrationTestSuite) Test_Integration_Receive_Heartbeat() {
	input, _ := hex.DecodeString(
		`adbccbda00000002000000000000000657534a542d580000000300000005322e322e3200000006306439623936`)
	expected := wsjtx.HeartbeatMessage{
		Id:        "WSJT-X",
		MaxSchema: 3,
		Version:   "2.2.2",
		Revision:  "0d9b96",
	}
	_, err := s.fake.SendMessage(input)
	s.Require().NoError(err)

	for {
		select {
		case msg := <-s.msgChan:
			switch msg.(type) {
			case wsjtx.HeartbeatMessage:
				actual := msg.(wsjtx.HeartbeatMessage)
				s.Require().Equal(expected, actual)
				return
			default:
				s.Failf("wrong message type", "expected type %T but got %T",
					wsjtx.HeartbeatMessage{}, msg)
				return
			}
		case err := <-s.errChan:
			s.Require().NoError(err)
			return
		case <-time.After(50 * time.Millisecond):
			s.Fail("timeout")
			return
		}
	}
}

func (s *integrationTestSuite) Test_Integration_Receive_Close() {
	input, _ := hex.DecodeString(
		`adbccbda00000002000000060000000657534a542d58`)
	expected := wsjtx.CloseMessage{
		Id: "WSJT-X",
	}
	_, err := s.fake.SendMessage(input)
	s.Require().NoError(err)

	for {
		select {
		case msg := <-s.msgChan:
			switch msg.(type) {
			case wsjtx.CloseMessage:
				actual := msg.(wsjtx.CloseMessage)
				s.Require().Equal(expected, actual)
				return
			default:
				s.Failf("wrong message type", "expected type %T but got %T",
					wsjtx.CloseMessage{}, msg)
				return
			}
		case err := <-s.errChan:
			s.Require().NoError(err)
			return
		case <-time.After(50 * time.Millisecond):
			s.Fail("timeout")
			return
		}
	}
}

func (s *integrationTestSuite) Test_Integration_Receive_Clear() {
	input, _ := hex.DecodeString(
		`adbccbda00000002000000030000000657534a542d58`)
	expected := wsjtx.ClearMessage{
		Id: "WSJT-X",
	}
	_, err := s.fake.SendMessage(input)
	s.Require().NoError(err)

	for {
		select {
		case msg := <-s.msgChan:
			switch msg.(type) {
			case wsjtx.ClearMessage:
				actual := msg.(wsjtx.ClearMessage)
				s.Require().Equal(expected, actual)
				return
			default:
				s.Failf("wrong message type", "expected type %T but got %T",
					wsjtx.ClearMessage{}, msg)
				return
			}
		case err := <-s.errChan:
			s.Require().NoError(err)
			return
		case <-time.After(50 * time.Millisecond):
			s.Fail("timeout")
			return
		}
	}
}
