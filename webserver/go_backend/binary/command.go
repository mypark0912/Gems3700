package binary

import "encoding/binary"

const (
	CmdClear   = 0
	CmdReboot  = 1
	CmdCapture = 2
)

const (
	ItemDemand   = 0
	ItemMaxMin   = 1
	ItemEnergy   = 2
	ItemAlarm    = 3
	ItemEvent    = 4
	ItemReboot   = 5
	ItemRunHour  = 6
	ItemWaveform = 7
	ItemAll      = 8
)

type Command struct {
	Type int32
	Cmd  int32
	Item int32
}

func (c Command) Encode() []byte {
	buf := make([]byte, 12)
	binary.LittleEndian.PutUint32(buf[0:4], uint32(c.Type))
	binary.LittleEndian.PutUint32(buf[4:8], uint32(c.Cmd))
	binary.LittleEndian.PutUint32(buf[8:12], uint32(c.Item))
	return buf
}

func DecodeCommand(data []byte) Command {
	if len(data) < 12 {
		return Command{}
	}
	return Command{
		Type: int32(binary.LittleEndian.Uint32(data[0:4])),
		Cmd:  int32(binary.LittleEndian.Uint32(data[4:8])),
		Item: int32(binary.LittleEndian.Uint32(data[8:12])),
	}
}

var ItemNames = map[int]string{
	0: "Demand Data",
	1: "Max/Min Data",
	2: "Energy Data",
	3: "Alarm Count",
	4: "Event Count",
	5: "Reboot",
	6: "RunHour Data",
	7: "Waveform",
	8: "All Data",
}
