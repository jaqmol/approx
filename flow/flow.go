package flow

import (
	"github.com/jaqmol/approx/conf"
)

// Flow ...
type Flow struct {
	MainItem *ProcItem
	// structure [][]Item
}

// NewFlow ...
func NewFlow(form *conf.Formation) *Flow {
	mainItem := NewProcItem(form.MainConf)
	procItemForName := make(map[string]*ProcItem)
	procItemForName[mainItem.Conf.Name()] = mainItem

	for name, c := range form.PrivateConfs {
		procItemForName[name] = NewProcItem(c)
	}

	for _, piA := range procItemForName {
		for _, piB := range procItemForName {
			if piA.IsConnected(piB.Conf) {
				continue
			}
			if piA.ShouldConnect(piB.Conf) {
				piA.Connect(piB)
				if piA == mainItem {
					for _, conn := range piB.Next {
						if conn.To == mainItem {
							conn.IsLoopBack = true
						}
					}
				}
			}
		}
	}

	return &Flow{
		MainItem: mainItem,
	}
}

// Iterate ...
func (f *Flow) Iterate(ito func(row []*ProcItem)) {
	currentProcs := NewProcSet(f.MainItem)
	loopBackProcs := NewProcSet(f.MainItem)
	for currentProcs.Len() > 0 {
		ito(currentProcs.Data())
		nextProcs := NewProcSet(nil)
		for _, proc := range currentProcs.Data() {
			for _, conn := range proc.Next {
				if conn.To == f.MainItem {
					loopBackProcs.Add(conn.To)
				} else {
					nextProcs.Add(conn.To)
				}
			}
		}
		currentProcs = nextProcs
	}
	if loopBackProcs.Len() > 0 {
		ito(loopBackProcs.Data())
	}
}
