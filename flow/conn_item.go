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

// // Type ...
// func (ci *ConnItem) Type() ItemType {
// 	return ItemTypeConn
// }

// // PreviousCount ...
// func (ci *ConnItem) PreviousCount() int {
// 	return 1
// }

// // Previous ...
// func (ci *ConnItem) Previous(int) Item {
// 	return ci.from
// }

// // NextCount ...
// func (ci *ConnItem) NextCount() int {
// 	return 1
// }

// // Next ...
// func (ci *ConnItem) Next(int) Item {
// 	return ci.to
// }

// // Visit ...
// func (ci *ConnItem) Visit(ito func(Item)) {
// 	ito(ci)
// 	if !ci.isLoopBack && ci.to != nil {
// 		ci.to.Visit(ito)
// 	}
// }
