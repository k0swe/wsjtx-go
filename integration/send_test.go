package integration

import (
	"time"

	"github.com/k0swe/wsjtx-go/v3"
)

func (s *integrationTestSuite) TestSendHeartbeat() {
	s.primeConnection()

	msg := wsjtx.HeartbeatMessage{
		Id:        "WSJT-X",
		MaxSchema: 3,
		Version:   "2.2.2",
		Revision:  "0d9b96"}
	want := decode(`adbccbda00000002000000000000000657534a542d580000000300000005322e322e3200000006306439623936`)

	s.T().Log("test case sending heartbeat struct")
	err := s.server.Heartbeat(msg)
	s.Require().NoError(err)
	s.waitForReceiveAndCheck(want)
}

func (s *integrationTestSuite) waitForReceiveAndCheck(want []byte) {
	for {
		select {
		case got := <-s.fake.ReceiveChan:
			s.T().Log("test case got receive bytes back from fake")
			s.Require().Equal(want, got)
			return
		case <-time.After(50 * time.Millisecond):
			s.Fail("timeout")
			return
		}
	}
}

func (s *integrationTestSuite) primeConnection() {
	// Because this is UDP, the server doesn't have an address for WSJTX until WSJTX has sent the
	// server a message.
	clearMsg := decode(`adbccbda00000002000000030000000657534a542d58`)
	_, err := s.fake.SendMessage(clearMsg)
	s.Require().NoError(err)
	<-s.msgChan
	s.T().Log("connection is primed for a send test")
}
