package wsjtx

import (
	"time"
)

/*
The heartbeat  message shall be  sent on a periodic  basis every
15   seconds.  This
message is intended to be used by servers to detect the presence
of a  client and also  the unexpected disappearance of  a client
and  by clients  to learn  the schema  negotiated by  the server
after it receives  the initial heartbeat message  from a client.

Out/In.

https://sourceforge.net/p/wsjt/wsjtx/ci/wsjtx-2.5.2/tree/Network/NetworkMessage.hpp#l110
*/
type HeartbeatMessage struct {
	Id        string `json:"id"`
	MaxSchema uint32 `json:"maxSchemaVersion"`
	Version   string `json:"version"`
	Revision  string `json:"revision"`
}

const heartbeatNum = 0

/*
WSJT-X  sends this  status message  when various  internal state
changes to allow the server to  track the relevant state of each
client without the need for  polling commands.

Out only.

https://sourceforge.net/p/wsjt/wsjtx/ci/wsjtx-2.5.2/tree/Network/NetworkMessage.hpp#l141
*/
type StatusMessage struct {
	Id                   string `json:"id"`
	DialFrequency        uint64 `json:"dialFrequency"`
	Mode                 string `json:"mode"`
	DxCall               string `json:"dxCall"`
	Report               string `json:"report"`
	TxMode               string `json:"txMode"`
	TxEnabled            bool   `json:"txEnabled"`
	Transmitting         bool   `json:"transmitting"`
	Decoding             bool   `json:"decoding"`
	RxDF                 uint32 `json:"rxDeltaFreq"`
	TxDF                 uint32 `json:"txDeltaFreq"`
	DeCall               string `json:"deCall"`
	DeGrid               string `json:"deGrid"`
	DxGrid               string `json:"dxGrid"`
	TxWatchdog           bool   `json:"txWatchdog"`
	SubMode              string `json:"submode"`
	FastMode             bool   `json:"fastMode"`
	SpecialOperationMode uint8  `json:"specialMode"`
	FrequencyTolerance   uint32 `json:"frequencyTolerance"`
	TRPeriod             uint32 `json:"txRxPeriod"`
	ConfigurationName    string `json:"configName"`
	TxMessage            string `json:"txMessage"`
}

const statusNum = 1

/*
The decode message is sent when  a new decode is completed, in
this case the 'New' field is true. It is also used in response
to  a "Replay"  message where  each  old decode  in the  "Band
activity" window, that  has not been erased, is  sent in order
as a one of these messages  with the 'New' field set to false.

Out only.

https://sourceforge.net/p/wsjt/wsjtx/ci/wsjtx-2.5.2/tree/Network/NetworkMessage.hpp#l208
*/
type DecodeMessage struct {
	Id               string  `json:"id"`
	New              bool    `json:"new"`
	Time             uint32  `json:"time"`
	Snr              int32   `json:"snr"`
	DeltaTimeSec     float64 `json:"deltaTime"`
	DeltaFrequencyHz uint32  `json:"deltaFrequency"`
	Mode             string  `json:"mode"`
	Message          string  `json:"message"`
	LowConfidence    bool    `json:"lowConfidence"`
	OffAir           bool    `json:"offAir"`
}

const decodeNum = 2

/*
This message is  send when all prior "Decode"  messages in the
"Band Activity"  window have been discarded  and therefore are
no long available for actioning  with a "Reply" message.

The Window  argument  can be  one  of the  following values:

	0  - clear the "Band Activity" window (default)
	1  - clear the "Rx Frequency" window
	2  - clear both "Band Activity" and "Rx Frequency" windows

Out/In.

https://sourceforge.net/p/wsjt/wsjtx/ci/wsjtx-2.5.2/tree/Network/NetworkMessage.hpp#l234
*/
type ClearMessage struct {
	Id     string `json:"id"`
	Window uint8  `json:"window"` // In only
}

const clearNum = 3

/*
In order for a server  to provide a useful cooperative service
to WSJT-X it  is possible for it to initiate  a QSO by sending
this message to a client. WSJT-X filters this message and only
acts upon it  if the message exactly describes  a prior decode
and that decode  is a CQ or QRZ message.   The action taken is
exactly equivalent to the user  double clicking the message in
the "Band activity" window.

In only.

https://sourceforge.net/p/wsjt/wsjtx/ci/wsjtx-2.5.2/tree/Network/NetworkMessage.hpp#l255
*/
type ReplyMessage struct {
	Id               string  `json:"id"`
	Time             uint32  `json:"time"`
	Snr              int32   `json:"snr"`
	DeltaTimeSec     float64 `json:"deltaTime"`
	DeltaFrequencyHz uint32  `json:"deltaFrequency"`
	Mode             string  `json:"mode"`
	Message          string  `json:"message"`
	LowConfidence    bool    `json:"lowConfidence"`
	Modifiers        uint8   `json:"modifiers"`
}

const replyNum = 4

/*
The QSO logged message is sent when the WSJT-X user accepts the "Log  QSO" dialog by clicking
the "OK" button.

Out only.

https://sourceforge.net/p/wsjt/wsjtx/ci/wsjtx-2.5.2/tree/Network/NetworkMessage.hpp#l295
*/
type QsoLoggedMessage struct {
	Id                  string    `json:"id"`
	DateTimeOff         time.Time `json:"dateTimeOff"`
	DxCall              string    `json:"dxCall"`
	DxGrid              string    `json:"dxGrid"`
	TxFrequency         uint64    `json:"txFrequency"`
	Mode                string    `json:"mode"`
	ReportSent          string    `json:"reportSent"`
	ReportReceived      string    `json:"reportReceived"`
	TxPower             string    `json:"txPower"`
	Comments            string    `json:"comments"`
	Name                string    `json:"name"`
	DateTimeOn          time.Time `json:"dateTimeOn"`
	OperatorCall        string    `json:"operatorCall"`
	MyCall              string    `json:"myCall"`
	MyGrid              string    `json:"myGrid"`
	ExchangeSent        string    `json:"exchangeSent"`
	ExchangeReceived    string    `json:"exchangeReceived"`
	ADIFPropagationMode string    `json:"propagationMode"`
}

const qsoLoggedNum = 5

/*
Close is  sent by  a client immediately  prior to  it shutting
down gracefully.

Out/In.

https://sourceforge.net/p/wsjt/wsjtx/ci/wsjtx-2.5.2/tree/Network/NetworkMessage.hpp#l320
*/
type CloseMessage struct {
	Id string `json:"id"`
}

const closeNum = 6

/*
When a server starts it may  be useful for it to determine the
state  of preexisting  clients. Sending  this message  to each
client as it is discovered  will cause that client (WSJT-X) to
send a "Decode" message for each decode currently in its "Band
activity"  window.

In only.

https://sourceforge.net/p/wsjt/wsjtx/ci/wsjtx-2.5.2/tree/Network/NetworkMessage.hpp#l328
*/
type ReplayMessage struct {
	Id string `json:"id"`
}

const replayNum = 7

/*
The server may stop a client from transmitting messages either
immediately or at  the end of the  current transmission period
using this message.

In only.

https://sourceforge.net/p/wsjt/wsjtx/ci/wsjtx-2.5.2/tree/Network/NetworkMessage.hpp#l343
*/
type HaltTxMessage struct {
	Id         string `json:"id"`
	AutoTxOnly bool   `json:"autoTxOnly"`
}

const haltTxNum = 8

/*
This message  allows the server  to set the current  free text
message content. Sending this  message with a non-empty "Text"
field is equivalent to typing  a new message (old contents are
discarded) in to  the WSJT-X free text message  field or "Tx5"
field (both  are updated) and if  the "Send" flag is  set then
clicking the "Now" radio button for the "Tx5" field if tab one
is current or clicking the "Free  msg" radio button if tab two
is current.

In only.

https://sourceforge.net/p/wsjt/wsjtx/ci/wsjtx-2.5.2/tree/Network/NetworkMessage.hpp#l352
*/
type FreeTextMessage struct {
	Id   string `json:"id"`
	Text string `json:"text"`
	Send bool   `json:"send"`
}

const freeTextNum = 9

/*
The decode message is sent when  a new decode is completed, in
this case the 'New' field is true.

Out only.

https://sourceforge.net/p/wsjt/wsjtx/ci/wsjtx-2.5.2/tree/Network/NetworkMessage.hpp#l383
*/
type WSPRDecodeMessage struct {
	Id        string  `json:"id"`
	New       bool    `json:"new"`
	Time      uint32  `json:"time"`
	Snr       int32   `json:"snr"`
	DeltaTime float64 `json:"deltaTime"`
	Frequency uint64  `json:"frequency"`
	Drift     int32   `json:"drift"`
	Callsign  string  `json:"callsign"`
	Grid      string  `json:"grid"`
	Power     int32   `json:"power"`
	OffAir    bool    `json:"offAir"`
}

const wsprDecodeNum = 10

/*
This  message allows  the server  to set  the current  current
geographical location  of operation. The supplied  location is
not persistent but  is used as a  session lifetime replacement
loction that overrides the Maidenhead  grid locater set in the
application  settings.

In only.

https://sourceforge.net/p/wsjt/wsjtx/ci/wsjtx-2.5.2/tree/Network/NetworkMessage.hpp#l406
*/
type LocationMessage struct {
	Id       string `json:"id"`
	Location string `json:"location"`
}

const locationNum = 11

/*
The  logged ADIF  message is  sent to  the server(s)  when the
WSJT-X user accepts the "Log  QSO" dialog by clicking the "OK"
button.

Out only.

https://sourceforge.net/p/wsjt/wsjtx/ci/wsjtx-2.5.2/tree/Network/NetworkMessage.hpp#l423
*/
type LoggedAdifMessage struct {
	Id   string `json:"id"`
	Adif string `json:"adif"`
}

const loggedAdifNum = 12

/*
The server  may send  this message at  any time.   The message
specifies  the background  and foreground  color that  will be
used  to  highlight  the  specified callsign  in  the  decoded
messages  printed  in the  Band  Activity  panel. To clear
and  cancel  highlighting send  an  invalid  QColor value  for
either or both  of the background and  foreground fields.

In only.

https://sourceforge.net/p/wsjt/wsjtx/ci/wsjtx-2.5.2/tree/Network/NetworkMessage.hpp#l444
*/
type HighlightCallsignMessage struct {
	Id              string `json:"id"`
	Callsign        string `json:"callsign"`
	BackgroundColor string `json:"backgroundColor"`
	ForegroundColor string `json:"foregroundColor"`
	HighlightLast   bool   `json:"highlightLast"`
	// This field is not part of the WSJT-X message and is specific to the golang library. It is a
	// necessary addition to be able to reset the highlighting. QT's color has a sentinel value in
	// QColor to signal an "invalid" color; golang image/color doesn't have that, so we add this
	// field. If this is true, BackgroundColor and ForegroundColor become "invalid" colors.
	Reset bool `json:"reset"`
}

const highlightCallsignNum = 13

/*
The server  may send  this message at  any time.   The message
specifies the name of the  configuration to switch to. The new
configuration must exist.

In only.

https://sourceforge.net/p/wsjt/wsjtx/ci/wsjtx-2.5.2/tree/Network/NetworkMessage.hpp#l445
*/
type SwitchConfigurationMessage struct {
	Id                string `json:"id"`
	ConfigurationName string `json:"configurationName"`
}

const switchConfigurationNum = 14

/*
The server  may send  this message at  any time.   The message
specifies  various  configuration  options.  For  utf8  string
fields an empty value implies no change, for the quint32 Rx DF
and  Frequency  Tolerance  fields the  maximum  quint32  value
implies  no change.   Invalid or  unrecognized values  will be
silently ignored.

In only.

https://sourceforge.net/p/wsjt/wsjtx/ci/wsjtx-2.5.2/tree/Network/NetworkMessage.hpp#l479
*/
type ConfigureMessage struct {
	Id                 string `json:"id"`
	Mode               string `json:"mode"`
	FrequencyTolerance uint32 `json:"frequencyTolerance"`
	Submode            string `json:"submode"`
	FastMode           bool   `json:"fastMode"`
	TRPeriod           uint32 `json:"trPeriod"`
	RxDF               uint32 `json:"rxDF"`
	DXCall             string `json:"dxCall"`
	DXGrid             string `json:"dxGrid"`
	GenerateMessages   bool   `json:"generateMessages"`
}

const configureNum = 15
