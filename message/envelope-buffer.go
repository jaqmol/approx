package message

import (
	"bytes"
	"log"
	"regexp"
	"strconv"

	"github.com/jaqmol/approx/utils"
)

// EnvelopeBuffer ...
type EnvelopeBuffer struct {
	buffer             []byte
	envs               chan *Envelope
	lengthPrefixRegex  *regexp.Regexp
	invalidSyntaxRegex *regexp.Regexp
}

// NewEnvelopeBuffer ...
func NewEnvelopeBuffer(chunks <-chan []byte) *EnvelopeBuffer {
	cc := EnvelopeBuffer{
		envs:               make(chan *Envelope),
		lengthPrefixRegex:  regexp.MustCompile("^\\d+:"),
		invalidSyntaxRegex: regexp.MustCompile("invalid syntax$"),
	}
	go cc.start(chunks)
	return &cc
}

func (c *EnvelopeBuffer) start(chunks <-chan []byte) {
	for chunk := range chunks {
		c.buffer = append(c.buffer, chunk...)
		msgLen, envBytes, ok := c.extractEnvelopeBytes()
		if ok {
			env := Envelope{
				MessageLength: msgLen,
				Bytes:         envBytes,
			}
			c.envs <- &env
		}
	}
}

// Envelopes ...
func (c *EnvelopeBuffer) Envelopes() <-chan *Envelope {
	return c.envs
}

func (c *EnvelopeBuffer) extractEnvelopeBytes() (msgLen int, envBytes []byte, ok bool) {
	msgLen = 0
	envBytes = nil
	ok = false
	if len(c.buffer) == 0 {
		return
	}
	envLen, msgLen := c.envelopeLength()
	if envLen == -1 {
		return
	}
	if len(c.buffer) < envLen {
		return
	}
	envBytes = c.buffer[0:envLen]
	c.buffer = c.buffer[envLen:len(c.buffer)]
	ok = true
	return
}

func (c *EnvelopeBuffer) envelopeLength() (envLen int, msgLen int) {
	colonIndex := bytes.IndexByte(c.buffer, ':')
	if colonIndex == -1 {
		return -1, -1
	}
	msgLenBuff := c.buffer[0:colonIndex]
	uMsgLen, err := strconv.ParseUint(string(msgLenBuff), 10, 64)
	if err != nil {
		if c.invalidSyntaxRegex.MatchString(err.Error()) && !c.lengthPrefixRegex.Match(c.buffer) {
			log.Fatalf("Invalid envelope synthax: %v\n", string(utils.Truncated(c.buffer, 80)))
		} else {
			log.Fatal(err.Error())
		}
	}
	if uMsgLen == 0 {
		panic("Length-prefix must always be > 0")
	}
	msgLen = int(uMsgLen)
	envLen = len(msgLenBuff) + 1 + msgLen
	return
}
