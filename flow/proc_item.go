package flow

import (
	"github.com/jaqmol/approx/conf"
)

// ProcItem ...
type ProcItem struct {
	Conf          conf.Conf
	Prev          []*ConnItem
	Next          []*ConnItem
	prevProcNames map[string]bool
	nextProcNames map[string]bool
}

// NewProcItem ...
func NewProcItem(c conf.Conf) *ProcItem {
	return &ProcItem{
		Conf:          c,
		Prev:          make([]*ConnItem, 0),
		Next:          make([]*ConnItem, 0),
		prevProcNames: make(map[string]bool),
		nextProcNames: make(map[string]bool),
	}
}

// IsConnected ...
func (pi *ProcItem) IsConnected(c conf.Conf) bool {
	confName := c.Name()
	ok1, ok2 := pi.nextProcNames[confName]
	if ok2 {
		return ok1
	}
	ok1, ok2 = pi.prevProcNames[confName]
	if ok2 {
		return ok1
	}
	return false
}

// ShouldConnect ...
func (pi *ProcItem) ShouldConnect(c conf.Conf) bool {
	confName := c.Name()
	connectedNames := append([]string{}, pi.Conf.Outputs()...)
	connectedNames = append(connectedNames, pi.Conf.Inputs()...)
	for _, connName := range connectedNames {
		if connName == confName {
			return true
		}
	}
	return false
}

// Connect ...
func (pi *ProcItem) Connect(toPI *ProcItem) *ProcItem {
	procItemA := pi
	procItemB := toPI
	procItemBName := procItemB.Conf.Name()
	connections := make([]*ConnItem, 0)
	for _, aOutputName := range procItemA.Conf.Outputs() {
		if aOutputName == procItemBName {
			pi.nextProcNames[procItemBName] = true
			conn := NewConnItem(procItemA, procItemB)
			connections = append(connections, conn)
			procItemA.Next = append(procItemA.Next, conn)
			procItemB.Prev = append(procItemB.Prev, conn)
		}
	}
	for _, aInputName := range procItemA.Conf.Inputs() {
		if aInputName == procItemBName {
			ok1, ok2 := pi.prevProcNames[procItemBName]
			if !(ok2 && ok1) {
				pi.prevProcNames[procItemBName] = true
				conn := NewConnItem(procItemB, procItemA)
				connections = append(connections, conn)
				procItemB.Next = append(procItemB.Next, conn)
				procItemA.Prev = append(procItemA.Prev, conn)
			}
		}
	}
	return toPI
}
