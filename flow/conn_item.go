package flow

import (
	"hash/fnv"
)

// ConnItem ...
type ConnItem struct {
	From       *ProcItem
	To         *ProcItem
	IsLoopBack bool
	Hash       uint32
}

// NewConnItem ...
func NewConnItem(from *ProcItem, to *ProcItem) *ConnItem {
	h := fnv.New32a()
	h.Write([]byte(from.Conf.Name()))
	h.Write([]byte(to.Conf.Name()))

	return &ConnItem{
		From: from,
		To:   to,
		Hash: h.Sum32(),
	}
}
