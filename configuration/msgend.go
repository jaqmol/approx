package configuration

// MsgEnd ...
const MsgEnd = "#APRX:MSGEND#"

// MsgEndBytes ...
var MsgEndBytes []byte

func init() {
	MsgEndBytes = []byte(MsgEnd)
}
