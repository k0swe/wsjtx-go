package wsjtx

import (
	"encoding/hex"
	"reflect"
	"testing"
	"time"
)

type parseArgs struct {
	buffer []byte
	length int
}

func TestParseMessage(t *testing.T) {

	tests := []struct {
		name string
		args parseArgs
		want interface{}
	}{
		{
			name: "Parse Heartbeat",
			args: argsFrom(`adbccbda00000002000000000000000657534a542d580000000300000005322e322e3200000006306439623936`),
			want: HeartbeatMessage{
				Id:        "WSJT-X",
				MaxSchema: 3,
				Version:   "2.2.2",
				Revision:  "0d9b96",
			},
		},
		{
			name: "Parse Status",
			args: argsFrom(`adbccbda00000002000000010000000657534a542d5800000000006bf0d000000003465438ffffffff000000032d313500000003465438000000000003730000079e000000054b3053574500000006444d37394c56ffffffff00ffffffff0000ffffffffffffffff0000000744656661756c74`),
			want: StatusMessage{
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
			},
		},
		{
			name: "Parse Decode",
			args: argsFrom(`adbccbda00000002000000020000000657534a542d58010259baf8fffffffb3fc99999a000000000000516000000017e0000000e4a4132454a50204e3442502037330000`),
			want: DecodeMessage{
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
			},
		},
		{
			name: "Parse Clear",
			args: argsFrom(`adbccbda00000002000000030000000657534a542d58`),
			want: ClearMessage{
				Id: "WSJT-X",
			},
		},
		{
			name: "Parse QSO Logged",
			args: argsFrom(`adbccbda00000002000000050000000657534a542d5800000000002586110277ac48010000000454335354000000044a4b373300000000006bf86e00000003465438000000022d33000000022d37000000013500000007436f6d6d656e74000000034a6f6500000000002586110276c1e801000000055433535452000000054b3053574500000006444d37394c56000000023142000000023144`),
			want: QsoLoggedMessage{
				Id:               "WSJT-X",
				DateTimeOff:      parseTime("2020-10-30 11:29:57 +0000 UTC"),
				DxCall:           "T3ST",
				DxGrid:           "JK73",
				TxFrequency:      7075950,
				Mode:             "FT8",
				ReportSent:       "-3",
				ReportReceived:   "-7",
				TxPower:          "5",
				Comments:         "Comment",
				Name:             "Joe",
				DateTimeOn:       parseTime("2020-10-30 11:28:57 +0000 UTC"),
				OperatorCall:     "T3STR",
				MyCall:           "K0SWE",
				MyGrid:           "DM79LV",
				ExchangeSent:     "1B",
				ExchangeReceived: "1D",
			},
		}, {
			name: "Parse Close",
			args: argsFrom(`adbccbda00000002000000060000000657534a542d58`),
			want: CloseMessage{
				Id: "WSJT-X",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseMessage(tt.args.buffer, tt.args.length); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func parseTime(str string) time.Time {
	ret, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", str)
	return ret
}

func argsFrom(str string) parseArgs {
	bytes, _ := hex.DecodeString(str)
	return parseArgs{
		buffer: bytes,
		length: len(bytes),
	}
}
