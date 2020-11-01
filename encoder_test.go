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
	wantBin, _ := hex.DecodeString("adbccbda00000002000000000000000657534a542d580000000300000005322e322e3200000006306439623936")
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
			want:    wantBin,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encodeHeartbeat(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("encodeClear() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("encodeClear() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_encodeClear(t *testing.T) {
	type args struct {
		msg ClearMessage
	}
	wantBin, _ := hex.DecodeString("adbccbda00000002000000030000000657534a542d5802")
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "encodeClear",
			args:    args{msg: ClearMessage{"WSJT-X", 2}},
			want:    wantBin,
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
				t.Errorf("encodeClear() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_encodeClose(t *testing.T) {
	type args struct {
		msg CloseMessage
	}
	wantBin, _ := hex.DecodeString("adbccbda00000002000000060000000657534a542d58")
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "encodeClose",
			args:    args{msg: CloseMessage{"WSJT-X"}},
			want:    wantBin,
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
				t.Errorf("encodeClose() got = %v, want %v", got, tt.want)
			}
		})
	}
}
