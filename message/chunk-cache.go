package message

import (
	"bytes"
	"strconv"
)

type chunkCache struct {
	cache        []byte
	msgBytesChan chan []byte
	// chunksChan   <-chan []byte
}

func newChunkCache(chunks <-chan []byte) *chunkCache {
	cc := chunkCache{}
	go cc.start(chunks)
	return &cc
}

func (c *chunkCache) start(chunks <-chan []byte) {
	for chunk := range chunks {
		c.cache = append(c.cache, chunk...)
		msgBytes, ok := c.extractMessageBytes()
		if ok {
			c.msgBytesChan <- msgBytes
		}
	}
}

func (c *chunkCache) messageBytes() <-chan []byte {
	return c.msgBytesChan
}

func (c *chunkCache) extractMessageBytes() ([]byte, bool) {
	if len(c.cache) == 0 {
		return nil, false
	}
	length, ok := c.extractMsgLength()
	if !ok {
		return nil, false
	}
	bytes := c.cache[0:length]
	c.cache = c.cache[length:len(c.cache)]
	return bytes, true
}

func (c *chunkCache) extractMsgLength() (int, bool) {
	colonIndex := bytes.IndexByte(c.cache, ':')
	if colonIndex == -1 {
		return 0, false
	}
	lenBuff := c.cache[0:colonIndex]
	c.cache = c.cache[colonIndex+1 : len(c.cache)]
	length, err := strconv.ParseInt(string(lenBuff), 10, 64)
	catch(err)
	if length <= 0 {
		return 0, false
	}
	return int(length), true
}
