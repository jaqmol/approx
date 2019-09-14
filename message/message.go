package message

import (
	"bytes"
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
	Body      []byte
}

// ParseMessage ...
func ParseMessage(byteSlice []byte) *Message {
	semicolonIndex := bytes.IndexRune(byteSlice, ';')
	if semicolonIndex == -1 {
		return nil
	}
	msgBytes := byteSlice[:semicolonIndex]
	msgString := string(msgBytes)
	comps := strings.Split(msgString, ",")
	if len(comps) != 8 {
		return nil
	}
	return &Message{
		ID:        comps[0],
		Role:      comps[1],
		Seq:       parseInt(comps[2]),
		IsEnd:     parseBool(comps[3]),
		Status:    parseInt(comps[4]),
		MediaType: comps[5],
		Encoding:  comps[6],
		Body:      byteSlice[semicolonIndex+1:],
	}
}

// ToBytes ...
func (m *Message) ToBytes() []byte {
	comps := []string{
		m.ID,
		m.Role,
		formatInt(m.Seq),
		strconv.FormatBool(m.IsEnd),
		formatInt(m.Status),
		m.MediaType,
		m.Encoding,
	}
	headerString := strings.Join(comps, ",")
	acc := []byte(headerString)
	acc = append(acc, []byte(";")...)
	acc = append(acc, m.Body...)
	acc = append(acc, []byte("\n")...)
	return acc
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
