package message

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// Message ...
type Message struct {
	ID        string
	Role      string
	Seq       int
	IsEnd     bool
	Status    int
	MediaType string
	Encoding  string
	Data      []byte
}

// ParseMessage expects arg byteSlice to be of format: <header>;<data>
func ParseMessage(byteSlice []byte) (*Message, error) {
	semicolonIndex := bytes.IndexRune(byteSlice, ';')
	if semicolonIndex == -1 {
		err := fmt.Errorf("Cannot parse message: %v", string(byteSlice))
		return nil, err
	}
	msgString := string(byteSlice[:semicolonIndex])
	comps := strings.Split(msgString, ",")
	return &Message{
		ID:        comps[0],
		Role:      comps[1],
		Seq:       parseInt(comps[2]),
		IsEnd:     parseBool(comps[3]),
		Status:    parseInt(comps[4]),
		MediaType: comps[5],
		Encoding:  comps[6],
		Data:      byteSlice[semicolonIndex+1:],
	}, nil
}

func (m *Message) headerBytes() []byte {
	comps := []string{
		m.ID,
		m.Role,
		formatInt(m.Seq),
		strconv.FormatBool(m.IsEnd),
		formatInt(m.Status),
		m.MediaType,
		m.Encoding,
	}
	str := strings.Join(comps, ",")
	bts := []byte(str)
	bts = append(bts, ';')
	return bts
}

// Envelope returns length-prefixed envelope of format: <byte-length-header-data>:<header>;<data>
func (m *Message) Envelope() *Envelope {
	return NewEnvelope(m.headerBytes(), m.Data)
}

func parseInt(comp string) int {
	value, err := strconv.ParseInt(comp, 10, 32)
	catch(err)
	return int(value)
}

func formatInt(value int) string {
	return strconv.FormatInt(int64(value), 10)
}

func parseBool(comp string) bool {
	value, err := strconv.ParseBool(comp)
	catch(err)
	return value
}

func catch(err error) {
	if err != nil {
		panic(err)
	}
}
