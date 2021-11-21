package integration

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/k0swe/wsjtx-go/v4"
)

type receiveCase struct {
	name string
	args []byte
	want receiveResult
}

type receiveResult struct {
	msg interface{}
	err error
}

func (s *integrationTestSuite) TestReceiveCases() {
	tests := []receiveCase{
		{
			name: "Heartbeat",
			args: decode(`adbccbda00000002000000000000000657534a542d580000000200000005322e322e3200000006306439623936`),
			want: receiveResult{wsjtx.HeartbeatMessage{
				Id:        "WSJT-X",
				MaxSchema: 2,
				Version:   "2.2.2",
				Revision:  "0d9b96",
			}, nil},
		},
		{
			name: "Status 2.2.2",
			args: decode(`adbccbda00000002000000010000000657534a542d5800000000006bf0d000000003465438ffffffff000000032d313500000003465438000000000003730000079e000000054b3053574500000006444d37394c56ffffffff00ffffffff0000ffffffffffffffff0000000744656661756c74`),
			want: receiveResult{wsjtx.StatusMessage{
				Id:                   "WSJT-X",
				DialFrequency:        7074000,
				Mode:                 "FT8",
				DxCall:               "",
				Report:               "-15",
				TxMode:               "FT8",
				TxEnabled:            false,
				Transmitting:         false,
				Decoding:             false,
				RxDF:                 883,
				TxDF:                 1950,
				DeCall:               "K0SWE",
				DeGrid:               "DM79LV",
				DxGrid:               "",
				TxWatchdog:           false,
				SubMode:              "",
				FastMode:             false,
				SpecialOperationMode: 0,
				FrequencyTolerance:   4294967295,
				TRPeriod:             4294967295,
				ConfigurationName:    "Default",
			}, wsjtx.ParseError},
		},
		{
			name: "Status 2.3.1",
			args: decode(`adbccbda00000002000000010000000657534a542d5800000000006bf0d000000003465438ffffffff000000032d313500000003465438000000000003730000079e000000054b3053574500000006444d37394c56ffffffff00ffffffff0000ffffffffffffffff0000000744656661756c7400000000`),
			want: receiveResult{wsjtx.StatusMessage{
				Id:                   "WSJT-X",
				DialFrequency:        7074000,
				Mode:                 "FT8",
				DxCall:               "",
				Report:               "-15",
				TxMode:               "FT8",
				TxEnabled:            false,
				Transmitting:         false,
				Decoding:             false,
				RxDF:                 883,
				TxDF:                 1950,
				DeCall:               "K0SWE",
				DeGrid:               "DM79LV",
				DxGrid:               "",
				TxWatchdog:           false,
				SubMode:              "",
				FastMode:             false,
				SpecialOperationMode: 0,
				FrequencyTolerance:   4294967295,
				TRPeriod:             4294967295,
				ConfigurationName:    "Default",
				TxMessage:            "",
			}, nil},
		},
		{
			name: "Decode",
			args: decode(`adbccbda00000002000000020000000657534a542d58010259baf8fffffffb3fc99999a000000000000516000000017e0000000e4a4132454a50204e3442502037330000`),
			want: receiveResult{wsjtx.DecodeMessage{
				Id:               "WSJT-X",
				New:              true,
				Time:             39435000,
				Snr:              -5,
				DeltaTimeSec:     0.20000000298023224,
				DeltaFrequencyHz: 1302,
				Mode:             "~",
				Message:          "JA2EJP N4BP 73",
				LowConfidence:    false,
				OffAir:           false,
			}, nil},
		},
		{
			name: "Clear",
			args: decode(`adbccbda00000002000000030000000657534a542d58`),
			want: receiveResult{wsjtx.ClearMessage{
				Id: "WSJT-X",
			}, nil},
		},
		{
			name: "QSO Logged",
			args: decode(`adbccbda00000002000000050000000657534a542d5800000000002586110277ac48010000000454335354000000044a4b373300000000006bf86e00000003465438000000022d33000000022d37000000013500000007436f6d6d656e74000000034a6f6500000000002586110276c1e801000000055433535452000000054b3053574500000006444d37394c5600000002314200000002314400000003494f4e`),
			want: receiveResult{wsjtx.QsoLoggedMessage{
				Id:                  "WSJT-X",
				DateTimeOff:         parseTime("2020-10-30 11:29:57 +0000 UTC"),
				DxCall:              "T3ST",
				DxGrid:              "JK73",
				TxFrequency:         7075950,
				Mode:                "FT8",
				ReportSent:          "-3",
				ReportReceived:      "-7",
				TxPower:             "5",
				Comments:            "Comment",
				Name:                "Joe",
				DateTimeOn:          parseTime("2020-10-30 11:28:57 +0000 UTC"),
				OperatorCall:        "T3STR",
				MyCall:              "K0SWE",
				MyGrid:              "DM79LV",
				ExchangeSent:        "1B",
				ExchangeReceived:    "1D",
				ADIFPropagationMode: "ION",
			}, nil},
		},
		{
			name: "Close",
			args: decode(`adbccbda00000002000000060000000657534a542d58`),
			want: receiveResult{wsjtx.CloseMessage{
				Id: "WSJT-X",
			}, nil},
		},
		{
			name: "WSPR Decode",
			args: decode(`adbccbda000000020000000a0000000657534a542d580102b5f840ffffffeebfe000000000000000000000006b6c7300000000000000054b3654475700000004434d39350000001700`),
			want: receiveResult{wsjtx.WSPRDecodeMessage{
				Id:        "WSJT-X",
				New:       true,
				Time:      45480000,
				Snr:       -18,
				DeltaTime: -0.5,
				Frequency: 7040115,
				Drift:     0,
				Callsign:  "K6TGW",
				Grid:      "CM95",
				Power:     23,
				OffAir:    false,
			}, nil},
		},
		{
			name: "Logged Adif",
			args: decode(`adbccbda000000020000000c0000000657534a542d580000015c0a3c616469665f7665723a353e332e312e300a3c70726f6772616d69643a363e57534a542d580a3c454f483e0a3c63616c6c3a343e54335354203c677269647371756172653a343e4a4b3733203c6d6f64653a333e465438203c7273745f73656e743a323e2d38203c7273745f726376643a323e2d39203c71736f5f646174653a383e3230323031303330203c74696d655f6f6e3a363e313230383136203c71736f5f646174655f6f66663a383e3230323031303330203c74696d655f6f66663a363e313230393136203c62616e643a333e34306d203c667265713a383e372e303735393530203c73746174696f6e5f63616c6c7369676e3a353e4b30535745203c6d795f677269647371756172653a363e444d37394c56203c74785f7077723a313e35203c636f6d6d656e743a373e436f6d6d656e74203c6e616d653a343e4a657373203c6f70657261746f723a353e5433535452203c454f523e`),
			want: receiveResult{wsjtx.LoggedAdifMessage{
				Id: "WSJT-X",
				Adif: `
<adif_ver:5>3.1.0
<programid:6>WSJT-X
<EOH>
<call:4>T3ST <gridsquare:4>JK73 <mode:3>FT8 <rst_sent:2>-8 <rst_rcvd:2>-9 <qso_date:8>20201030 <time_on:6>120816 <qso_date_off:8>20201030 <time_off:6>120916 <band:3>40m <freq:8>7.075950 <station_callsign:5>K0SWE <my_gridsquare:6>DM79LV <tx_pwr:1>5 <comment:7>Comment <name:4>Jess <operator:5>T3STR <EOR>`,
			}, nil},
		},
	}

	for _, tt := range tests {
		s.runReceiveTest(tt)
	}
}

func (s *integrationTestSuite) runReceiveTest(tt receiveCase) {
	s.T().Run(tt.name, func(t *testing.T) {
		msgPassed, errPassed := false, false
		if tt.want.err == nil {
			errPassed = true
		}
		_, err := s.fake.SendMessage(tt.args)
		s.Require().NoError(err)

		for {
			select {
			case msg := <-s.msgChan:
				s.Require().Equal(tt.want.msg, msg)
				msgPassed = true
				if msgPassed && errPassed {
					return
				}
			case chErr := <-s.errChan:
				s.Require().ErrorIs(chErr, tt.want.err)
				errPassed = true
				if msgPassed && errPassed {
					return
				}
			case <-time.After(50 * time.Millisecond):
				s.Fail("timeout")
				return
			}
		}
	})
}

func decode(s string) []byte {
	input, _ := hex.DecodeString(s)
	return input
}

func parseTime(str string) time.Time {
	ret, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", str)
	return ret
}
