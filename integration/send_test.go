package integration

import (
	"math"
	"time"

	"github.com/k0swe/wsjtx-go/v4"
)

func (s *integrationTestSuite) TestSendHeartbeat() {
	s.primeConnection()

	msg := wsjtx.HeartbeatMessage{
		Id:        "WSJT-X",
		MaxSchema: 3,
		Version:   "2.2.2",
		Revision:  "0d9b96"}
	want := decode(`adbccbda00000002000000000000000657534a542d580000000300000005322e322e3200000006306439623936`)

	s.T().Log("sending heartbeat struct")
	err := s.server.Heartbeat(msg)
	s.Require().NoError(err)
	s.waitForReceiveAndCheck(want)
}

func (s *integrationTestSuite) TestSendClear() {
	s.primeConnection()

	msg := wsjtx.ClearMessage{Id: "WSJT-X", Window: 2}
	want := decode(`adbccbda00000002000000030000000657534a542d5802`)

	s.T().Log("sending clear struct")
	err := s.server.Clear(msg)
	s.Require().NoError(err)
	s.waitForReceiveAndCheck(want)
}

func (s *integrationTestSuite) TestSendReply() {
	s.primeConnection()

	msg := wsjtx.ReplyMessage{
		Id:               "WSJT-X",
		Time:             1234,
		Snr:              -15,
		DeltaTimeSec:     0.5,
		DeltaFrequencyHz: 2345,
		Mode:             "FT8",
		Message:          "CQ K0SWE DM79",
		LowConfidence:    false,
		Modifiers:        0,
	}
	want := decode(`adbccbda00000002000000040000000657534a542d58000004d2fffffff13fe000000000000000000929000000034654380000000d4351204b3053574520444d37390000`)

	s.T().Log("sending reply struct")
	err := s.server.Reply(msg)
	s.Require().NoError(err)
	s.waitForReceiveAndCheck(want)
}

func (s *integrationTestSuite) TestSendClose() {
	s.primeConnection()

	msg := wsjtx.CloseMessage{Id: "WSJT-X"}
	want := decode(`adbccbda00000002000000060000000657534a542d58`)

	s.T().Log("sending close struct")
	err := s.server.Close(msg)
	s.Require().NoError(err)
	s.waitForReceiveAndCheck(want)
}

func (s *integrationTestSuite) TestSendReplay() {
	s.primeConnection()

	msg := wsjtx.ReplayMessage{Id: "WSJT-X"}
	want := decode(`adbccbda00000002000000070000000657534a542d58`)

	s.T().Log("sending replay struct")
	err := s.server.Replay(msg)
	s.Require().NoError(err)
	s.waitForReceiveAndCheck(want)
}

func (s *integrationTestSuite) TestSendHaltTx() {
	s.primeConnection()

	msg := wsjtx.HaltTxMessage{Id: "WSJT-X", AutoTxOnly: false}
	want := decode(`adbccbda00000002000000080000000657534a542d5800`)

	s.T().Log("sending haltTx struct")
	err := s.server.HaltTx(msg)
	s.Require().NoError(err)
	s.waitForReceiveAndCheck(want)
}

func (s *integrationTestSuite) TestSendFreeText() {
	s.primeConnection()

	msg := wsjtx.FreeTextMessage{
		Id:   "WSJT-X",
		Text: "J72IMS K0SWE R-15",
		Send: true,
	}
	want := decode(`adbccbda00000002000000090000000657534a542d58000000114a3732494d53204b3053574520522d313501`)

	s.T().Log("sending freeText struct")
	err := s.server.FreeText(msg)
	s.Require().NoError(err)
	s.waitForReceiveAndCheck(want)
}

func (s *integrationTestSuite) TestSendLocation() {
	s.primeConnection()

	msg := wsjtx.LocationMessage{Id: "WSJT-X", Location: "DM79jx"}
	want := decode(`adbccbda000000020000000b0000000657534a542d5800000006444d37396a78`)

	s.T().Log("sending location struct")
	err := s.server.Location(msg)
	s.Require().NoError(err)
	s.waitForReceiveAndCheck(want)
}

func (s *integrationTestSuite) TestSendHighlightCallsign() {
	s.primeConnection()

	msg := wsjtx.HighlightCallsignMessage{
		Id:              "WSJT-X",
		Callsign:        "KM4ACK",
		BackgroundColor: "red",
		ForegroundColor: "black",
		HighlightLast:   true,
		Reset:           false,
	}
	want := decode(`adbccbda000000020000000d0000000657534a542d58000000064b4d3441434b01ffffffff00000000000001ffff000000000000000001`)

	s.T().Log("sending highlightCallsign struct")
	err := s.server.HighlightCallsign(msg)
	s.Require().NoError(err)
	s.waitForReceiveAndCheck(want)
}

func (s *integrationTestSuite) TestSendSwitchConfiguration() {
	s.primeConnection()

	msg := wsjtx.SwitchConfigurationMessage{Id: "WSJT-X", ConfigurationName: "IC-7300"}
	want := decode(`adbccbda000000020000000e0000000657534a542d580000000749432d37333030`)

	s.T().Log("sending switchConfiguration struct")
	err := s.server.SwitchConfiguration(msg)
	s.Require().NoError(err)
	s.waitForReceiveAndCheck(want)
}

func (s *integrationTestSuite) TestSendConfigure() {
	s.primeConnection()

	msg := wsjtx.ConfigureMessage{
		Id:                 "WSJT-X",
		Mode:               "JT9",
		FrequencyTolerance: math.MaxUint32,
		Submode:            "",
		FastMode:           false,
		TRPeriod:           0,
		RxDF:               math.MaxUint32,
		DXCall:             "KI6NAZ",
		DXGrid:             "DM03",
		GenerateMessages:   false,
	}
	want := decode(`adbccbda000000020000000f0000000657534a542d58000000034a5439ffffffffffffffff0000000000ffffffff000000064b49364e415a00000004444d303300`)

	s.T().Log("sending configure struct")
	err := s.server.Configure(msg)
	s.Require().NoError(err)
	s.waitForReceiveAndCheck(want)
}

func (s *integrationTestSuite) waitForReceiveAndCheck(want []byte) {
	for {
		select {
		case got := <-s.fake.ReceiveChan:
			s.T().Log("got receive bytes back from fake")
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
