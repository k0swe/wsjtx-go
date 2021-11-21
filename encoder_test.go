package wsjtx

import (
	"encoding/hex"
	"reflect"
	"testing"
)

func Test_encodeHeartbeat(t *testing.T) {
	type args struct {
		msg HeartbeatMessage
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "encodeHeartbeat",
			args: args{msg: HeartbeatMessage{
				Id:        "WSJT-X",
				MaxSchema: 3,
				Version:   "2.2.2",
				Revision:  "0d9b96"}},
			want:    decodeHex("adbccbda00000002000000000000000657534a542d580000000300000005322e322e3200000006306439623936"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encodeHeartbeat(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("encodeHeartbeat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("encodeHeartbeat() got = %v, want %v",
					hex.EncodeToString(got), hex.EncodeToString(tt.want))
			}
		})
	}
}

func Test_encodeClear(t *testing.T) {
	type args struct {
		msg ClearMessage
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "encodeClear",
			args:    args{msg: ClearMessage{"WSJT-X", 2}},
			want:    decodeHex("adbccbda00000002000000030000000657534a542d5802"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encodeClear(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("encodeClear() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("encodeClear() got = %v, want %v",
					hex.EncodeToString(got), hex.EncodeToString(tt.want))
			}
		})
	}
}

func Test_encodeReply(t *testing.T) {
	type args struct {
		msg ReplyMessage
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "encodeReply",
			args: args{msg: ReplyMessage{
				Id:               "WSJT-X",
				Time:             39435000,
				Snr:              -5,
				DeltaTimeSec:     0.20000000298023224,
				DeltaFrequencyHz: 1302,
				Mode:             "~",
				Message:          "JA2EJP N4BP 73",
				LowConfidence:    false,
				Modifiers:        0,
			}},
			want:    decodeHex("adbccbda00000002000000040000000657534a542d580259baf8fffffffb3fc99999a000000000000516000000017e0000000e4a4132454a50204e3442502037330000"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encodeReply(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("encodeReply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("encodeReply() got = %v, want %v",
					hex.EncodeToString(got), hex.EncodeToString(tt.want))
			}
		})
	}
}

func Test_encodeClose(t *testing.T) {
	type args struct {
		msg CloseMessage
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "encodeClose",
			args:    args{msg: CloseMessage{"WSJT-X"}},
			want:    decodeHex("adbccbda00000002000000060000000657534a542d58"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encodeClose(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("encodeClose() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("encodeClose() got = %v, want %v",
					hex.EncodeToString(got), hex.EncodeToString(tt.want))
			}
		})
	}
}

func Test_encodeReplay(t *testing.T) {
	type args struct {
		msg ReplayMessage
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "encodeReplay",
			args:    args{msg: ReplayMessage{"WSJT-X"}},
			want:    decodeHex("adbccbda00000002000000070000000657534a542d58"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encodeReplay(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("encodeReplay() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("encodeReplay() got = %v, want %v",
					hex.EncodeToString(got), hex.EncodeToString(tt.want))
			}
		})
	}
}

func Test_encodeHaltTx(t *testing.T) {
	type args struct {
		msg HaltTxMessage
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "encodeHaltTx",
			args:    args{msg: HaltTxMessage{"WSJT-X", false}},
			want:    decodeHex("adbccbda00000002000000080000000657534a542d5800"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encodeHaltTx(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("encodeHaltTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("encodeHaltTx() got = %v, want %v",
					hex.EncodeToString(got), hex.EncodeToString(tt.want))
			}
		})
	}
}

func Test_encodeFreeText(t *testing.T) {
	type args struct {
		msg FreeTextMessage
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "encodeFreeText",
			args:    args{msg: FreeTextMessage{"WSJT-X", "😊 de 🇺🇸", true}},
			want:    decodeHex("adbccbda00000002000000090000000657534a542d5800000010f09f988a20646520f09f87baf09f87b801"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encodeFreeText(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("encodeFreeText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("encodeFreeText() got = %v, want %v",
					hex.EncodeToString(got), hex.EncodeToString(tt.want))
			}
		})
	}
}

func Test_encodeLocation(t *testing.T) {
	type args struct {
		msg LocationMessage
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "encodeLocation",
			args:    args{msg: LocationMessage{"WSJT-X", "DM79hu"}},
			want:    decodeHex("adbccbda000000020000000b0000000657534a542d5800000006444d37396875"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encodeLocation(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("encodeLocation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("encodeLocation() got = %v, want %v",
					hex.EncodeToString(got), hex.EncodeToString(tt.want))
			}
		})
	}
}

func Test_encodeHighlightCallsign(t *testing.T) {
	type args struct {
		msg HighlightCallsignMessage
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "encodeHighlightCallsign",
			args: args{msg: HighlightCallsignMessage{
				Id:              "WSJT-X",
				Callsign:        "KM4ACK",
				BackgroundColor: "#eb4034",
				ForegroundColor: "#25272e",
				HighlightLast:   true,
				Reset:           false,
			}},
			want:    decodeHex("adbccbda000000020000000d0000000657534a542d58000000064b4d3441434b01ffffebeb40403434000001ffff252527272e2e000001"),
			wantErr: false,
		},
		{
			name: "encodeHighlightCallsign_reset",
			args: args{msg: HighlightCallsignMessage{
				Id:              "WSJT-X",
				Callsign:        "KM4ACK",
				BackgroundColor: "white",
				ForegroundColor: "white",
				HighlightLast:   false,
				Reset:           true,
			}},
			want:    decodeHex("adbccbda000000020000000d0000000657534a542d58000000064b4d3441434b00ffffffffffffffff000000ffffffffffffffff000000"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encodeHighlightCallsign(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("encodeHighlightCallsign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("encodeHighlightCallsign() got = %v, want %v",
					hex.EncodeToString(got), hex.EncodeToString(tt.want))
			}
		})
	}
}

func Test_encodeSwitchConfiguration(t *testing.T) {
	type args struct {
		msg SwitchConfigurationMessage
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "encodeSwitchConfiguration",
			args:    args{msg: SwitchConfigurationMessage{"WSJT-X", "MyAlternateConfiguration"}},
			want:    decodeHex("adbccbda000000020000000e0000000657534a542d58000000184d79416c7465726e617465436f6e66696775726174696f6e"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encodeSwitchConfiguration(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("encodeSwitchConfiguration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("encodeSwitchConfiguration() got = %v, want %v",
					hex.EncodeToString(got), hex.EncodeToString(tt.want))
			}
		})
	}
}

func Test_encodeConfigure(t *testing.T) {
	type args struct {
		msg ConfigureMessage
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "encodeConfigure",
			args: args{msg: ConfigureMessage{
				Id:                 "WSJT-X",
				Mode:               "FT4",
				FrequencyTolerance: 35,
				Submode:            "",
				FastMode:           true,
				TRPeriod:           60,
				RxDF:               1000,
				DXCall:             "T3ST",
				DXGrid:             "JK73",
				GenerateMessages:   true,
			}},
			want:    decodeHex("adbccbda000000020000000f0000000657534a542d580000000346543400000023ffffffff010000003c000003e80000000454335354000000044a4b373301"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encodeConfigure(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("encodeConfigure() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("encodeConfigure() got = %v, want %v",
					hex.EncodeToString(got), hex.EncodeToString(tt.want))
			}
		})
	}
}

func decodeHex(str string) []byte {
	bits, _ := hex.DecodeString(str)
	return bits
}
