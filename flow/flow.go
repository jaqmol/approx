package flow

import (
	"github.com/jaqmol/approx/conf"
	"github.com/jaqmol/approx/errormsg"
)

// Init ...
func Init(errMsg *errormsg.ErrorMsg) (fl *Flow) {
	fo := conf.ReadFormation(errMsg)
	re := conf.NewReqEnv(fo)
	exitIfRequirementsAreMissing(errMsg, re)
	fl = NewFlow(fo)
	// hub, err = NewHub(re, fo)
	return
}

func exitIfRequirementsAreMissing(errMsg *errormsg.ErrorMsg, re *conf.ReqEnv) {
	allNames := make([]string, 0)
	for name, hasValue := range re.HasValuesForNames {
		allNames = append(allNames, name)
		if !hasValue {
			errMsg.LogFatal(nil, "Please provide environment variable: %v", name)
		}
	}
}

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

// IterateProcs ...
func (f *Flow) IterateProcs(ito func(row []*ProcItem)) {
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

// IterateConns ...
func (f *Flow) IterateConns(ito func(row []*ConnItem)) {
	currentConns := NewConnSet(f.MainItem.Next)
	loopBackConns := NewConnSet(f.MainItem.Prev)
	for currentConns.Len() > 0 {
		noLoopBackConns := NewConnSet(nil)
		for _, conn := range currentConns.Data() {
			if conn.To == f.MainItem {
				loopBackConns.Add(conn)
			} else {
				noLoopBackConns.Add(conn)
			}
		}
		ito(noLoopBackConns.Data())
		nextConns := NewConnSet(nil)
		for _, conn := range noLoopBackConns.Data() {
			nextConns.AddAll(conn.To.Next)
		}
		currentConns = nextConns
	}
	if loopBackConns.Len() > 0 {
		ito(loopBackConns.Data())
	}
}
